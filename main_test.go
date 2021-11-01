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
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	testCtx context.Context
	testDB  *pgxpool.Pool
)

func testDSN() string {
	params := map[string]string{
		"PGHOST":     "db",
		"PGDATABASE": "testdata",
		"PGUSER":     "testdata",
		"PGPORT":     "5432",
	}

	for k := range params {
		if v, ok := os.LookupEnv(k); ok {
			params[k] = v
		}
	}

	const dsnFmt = "host=%s dbname=%s user=%s port=%s"

	return fmt.Sprintf(dsnFmt, params["PGHOST"], params["PGDATABASE"], params["PGUSER"], params["PGPORT"])
}

func execQuerySlice(ctx context.Context, sqls []string) {
	for _, sql := range sqls {
		runWithCtxTimeout(ctx, 1*time.Second, func(c context.Context) {
			if _, err := testDB.Exec(c, sql); err != nil {
				log.Fatalf("testing.execQuerySlice: %v", err)
			}
		})
	}
}

//go:embed testdata/create.sql
var createTablesSQL string

//go:embed testdata/drop.sql
var dropTablesSQL string

func TestMain(m *testing.M) {
	var cancel context.CancelFunc
	testCtx, cancel = context.WithTimeout(context.Background(), 30*time.Second)

	runWithCtxTimeout(testCtx, 1*time.Second, func(c context.Context) {
		testDB = connectDB(c, testDSN())
	})

	execQuerySlice(testCtx, strings.SplitAfter(dropTablesSQL, ";"))
	execQuerySlice(testCtx, strings.SplitAfter(createTablesSQL, ";"))

	exit := m.Run()

	cancel()
	os.Exit(exit)
}

func Test_connectDB(t *testing.T) {
	tests := []struct {
		name    string
		dsn     string
		wantErr bool
	}{
		{
			"Failure",
			"foo",
			true,
		},
		{
			"Success",
			testDSN(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := func() (err error) {
				defer func() {
					err, _ = recover().(error)
				}()

				pool := connectDB(testCtx, tt.dsn)
				return pool.Ping(testCtx)
			}()

			if (err != nil) != tt.wantErr {
				t.Errorf("connectDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_acquireConn(t *testing.T) {
	ectx, cancel := context.WithCancel(testCtx)
	cancel()

	tests := []struct {
		name    string
		ctx     context.Context
		wantErr bool
	}{
		{
			"Context error",
			ectx,
			true,
		},
		{
			"Succes",
			testCtx,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := func() (err error) {
				defer func() { err, _ = recover().(error) }()
				conn := acquireConn(tt.ctx, testDB)
				return conn.Ping(testCtx)
			}()

			if (err != nil) != tt.wantErr {
				t.Errorf("acquireConn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_run(t *testing.T) {
	tests := []struct {
		name     string
		cf       string
		wantExit int
	}{
		{
			"Config error",
			"testdata/invalid.yml",
			1,
		},
		{
			"Success",
			"testdata/unit_test.yml",
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotExit := run(tt.cf); gotExit != tt.wantExit {
				t.Errorf("run() = %v, want %v", gotExit, tt.wantExit)
			}
		})
	}
}

type regressionColumn struct {
	Name string
	Data []interface{}
}

var recordRegressionData bool

func init() {
	flag.BoolVar(&recordRegressionData, "rec_reg_data", false, "Record regression data")
}

func Test_regression(t *testing.T) {
	exit := run("testdata/regression_test.yml")
	if exit != 0 {
		t.Fatal("regression test failed")
	}

	const sql = "select * from regression_tests;"

	var got []regressionColumn

	runWithCtxTimeout(testCtx, 5*time.Second, func(c context.Context) {
		rows, err := testDB.Query(c, sql)
		if err != nil {
			t.Fatalf("regression.Query: %v", err)
		}

		fds := rows.FieldDescriptions()
		got = make([]regressionColumn, len(fds))

		for i, fd := range fds {
			got[i].Name = string(fd.Name)
		}

		for rows.Next() {
			values, err := rows.Values()
			if err != nil {
				t.Fatal(err)
			}

			for i, v := range values {
				got[i].Data = append(got[i].Data, v)
			}
		}
	})

	const regressionDataFile = "testdata/regression.json"

	if recordRegressionData {
		file, err := os.Create(regressionDataFile)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		if err = enc.Encode(got); err != nil {
			t.Fatal(err)
		}

		if err = file.Sync(); err != nil {
			t.Fatal(err)
		}
		file.Close()
	}

	file, err := os.Open(regressionDataFile)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	var want []regressionColumn
	if err = dec.Decode(&want); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("regression got =\n%v\nwant\n%v", got, want)
	}
}
