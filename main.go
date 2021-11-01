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
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muhlemmer/pg_testdata/parse"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "conf", "pg_testdata.yml", "YAML config file with schema definitions")
}

func runWithCtxTimeout(ctx context.Context, d time.Duration, f func(context.Context)) {
	ctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()

	f(ctx)
}

func connectDB(ctx context.Context, dsn string) (pool *pgxpool.Pool) {
	runWithCtxTimeout(ctx, 5*time.Second, func(c context.Context) {
		var err error
		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			panic(fmt.Errorf("main.connectDB: %w", err))
		}
	})

	return pool
}

func acquireConn(ctx context.Context, pool *pgxpool.Pool) (conn *pgxpool.Conn) {
	runWithCtxTimeout(ctx, 5*time.Second, func(ctx context.Context) {
		var err error

		conn, err = pool.Acquire(ctx)
		if err != nil {
			panic(fmt.Errorf("main.acquire: %w", err))
		}
	})

	return conn
}

func run(cf string) (exit int) {
	defer func() {
		err, _ := recover().(error)
		if err != nil {
			log.Printf("FATAL ERROR: %v", err)
			exit = 1
		}

		var r runtime.Error
		if errors.As(err, &r) {
			panic(r)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	conf, err := parse.Load(cf)
	if err != nil {
		panic(err)
	}

	pool := connectDB(ctx, conf.DSN)

	for _, table := range conf.Tables {
		execInserts(ctx, pool, table)
	}

	return 0
}

func main() {
	flag.Parse()
	os.Exit(run(configFile))
}
