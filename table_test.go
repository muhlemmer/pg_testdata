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
	"testing"
	"time"

	"github.com/muhlemmer/pg_testdata/parse"
)

func Test_prepareInsert(t *testing.T) {
	ectx, cancel := context.WithCancel(testCtx)
	cancel()

	type args struct {
		ctx   context.Context
		table *parse.Table
	}
	tests := []struct {
		name        string
		args        args
		wantArgsLen int
		wantErr     bool
	}{
		{
			"InsertQuery error",
			args{
				testCtx,
				&parse.Table{
					Name:   "test_table",
					Amount: 5,
					MaxDuration: parse.TableDurations{
						Table: 10 * time.Second,
						Exec:  1 * time.Second,
					},
					Columns: []*parse.Column{
						{
							Name: "first",
							Type: "does-not-exist",
						},
					},
				},
			},
			0,
			true,
		},
		{
			"Context error",
			args{
				ectx,
				&parse.Table{
					Name:   "test_table",
					Amount: 5,
					MaxDuration: parse.TableDurations{
						Table: 10 * time.Second,
						Exec:  1 * time.Second,
					},
					Columns: []*parse.Column{
						{
							Name: "bool_col",
							Type: parse.BoolType,
							Generator: map[parse.ArgName]interface{}{
								parse.ProbabilityArg: 100,
							},
						},
					},
				},
			},
			0,
			true,
		},
		{
			"Context error",
			args{
				testCtx,
				&parse.Table{
					Name:   "unit_tests",
					Amount: 5,
					MaxDuration: parse.TableDurations{
						Table: 10 * time.Second,
						Exec:  1 * time.Second,
					},
					Columns: []*parse.Column{
						{
							Name: "bool_col",
							Type: parse.BoolType,
							Generator: map[parse.ArgName]interface{}{
								parse.ProbabilityArg: 100,
							},
						},
					},
				},
			},
			1,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := func() (err error) {
				defer func() { err, _ = recover().(error) }()

				conn := acquireConn(testCtx, testDB)
				defer conn.Release()

				_, gotArgs := prepareInsert(tt.args.ctx, conn, tt.args.table)
				if len(gotArgs) != tt.wantArgsLen {
					t.Errorf("prepareInsert() gotArgsLen = %d, want %d", gotArgs, tt.wantArgsLen)
				}

				return
			}()

			if (err != nil) != tt.wantErr {
				t.Errorf("prepareInsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_execInserts(t *testing.T) {
	tests := []struct {
		name    string
		table   *parse.Table
		wantErr bool
	}{
		{
			"Exec error",
			&parse.Table{
				Name:   "error_tests",
				Amount: 5,
				MaxDuration: parse.TableDurations{
					Table: 10 * time.Second,
					Exec:  1 * time.Second,
				},
				Columns: []*parse.Column{
					{
						Name: "bool_col",
						Type: parse.BoolType,
						Generator: map[parse.ArgName]interface{}{
							parse.ProbabilityArg: 100,
						},
					},
				},
			},
			true,
		},
		{
			"Succes",
			&parse.Table{
				Name:   "unit_tests",
				Amount: 5,
				MaxDuration: parse.TableDurations{
					Table: 10 * time.Second,
					Exec:  1 * time.Second,
				},
				Columns: []*parse.Column{
					{
						Name: "bool_col",
						Type: parse.BoolType,
						Generator: map[parse.ArgName]interface{}{
							parse.ProbabilityArg: 100,
						},
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := func() (err error) {
				defer func() { err, _ = recover().(error) }()
				execInserts(testCtx, testDB, tt.table)

				return
			}()

			if (err != nil) != tt.wantErr {
				t.Errorf("execInserts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})

	}
}
