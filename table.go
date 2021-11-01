/*
SPDX-License-Identifier: AGPL-3.0-only

pg_testdata is a test data generator for PostgreSQL.
Copyright (C) 2021  Tim Mohlmann

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muhlemmer/pg_testdata/parse"
)

func prepareInsert(ctx context.Context, conn *pgxpool.Conn, table *parse.Table) (sd *pgconn.StatementDescription, args []interface{}) {
	stmt, args, err := table.InsertQuery()
	if err != nil {
		panic(err)
	}

	runWithCtxTimeout(ctx, 5*time.Second, func(ctx context.Context) {
		sd, err = conn.Conn().Prepare(ctx, fmt.Sprintf("%s_insert", table.Name), stmt)
		if err != nil {
			panic(fmt.Errorf("main.prepareInsert: %w for table %q", err, table.Name))
		}
	})

	return sd, args
}

func execInserts(ctx context.Context, pool *pgxpool.Pool, table *parse.Table) {
	ctx, cancel := context.WithTimeout(ctx, table.MaxDuration.Table)
	defer cancel()

	conn := acquireConn(ctx, pool)
	sd, args := prepareInsert(ctx, conn, table)

	for i := 0; i < table.Amount; i++ {
		runWithCtxTimeout(ctx, table.MaxDuration.Exec, func(ctx context.Context) {
			if _, err := conn.Exec(ctx, sd.Name, args...); err != nil {
				panic(fmt.Errorf("main.execInsert: %w", err))
			}
		})
	}
}
