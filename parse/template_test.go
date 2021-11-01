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

package parse

import (
	"os"
	"testing"
)

func Test_lookupEnv(t *testing.T) {
	os.Setenv("test_key", "test_value")

	type args struct {
		key      string
		defValue string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Default",
			args{"foo", "bar"},
			"bar",
		},
		{
			"From env",
			args{"test_key", "bar"},
			"test_value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lookupEnv(tt.args.key, tt.args.defValue); got != tt.want {
				t.Errorf("lookupEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

const (
	yamlTemplateOutDefault = `dsn: dbname=testdata user=testdata host=db port=5432 connect_timeout=10
tables:
- name: unit_tests
  amount: 1000
  max_duration:
    table: 1m0s
    exec: 1s
  columns:
  - name: bool_col
    seed: 2
    nullprobability: 10
    type: bool
    generator:
      probability: 70.1
`
	yamlTemplateOutFromEnv = `dsn: dbname=foo user=bar host=db port=5432 connect_timeout=10
tables:
- name: unit_tests
  amount: 1000
  max_duration:
    table: 1m0s
    exec: 1s
  columns:
  - name: bool_col
    seed: 2
    nullprobability: 10
    type: bool
    generator:
      probability: 70.1
`
)

func Test_yamlTemplate(t *testing.T) {
	type invalidTmplData int

	type args struct {
		data     interface{}
		filename string
	}
	tests := []struct {
		name    string
		args    args
		env     map[string]string
		want    string
		wantErr bool
	}{
		{
			"File not found",
			args{
				nil,
				"/foo/bar/does/not/exist",
			},
			nil,
			"",
			true,
		},
		{
			"Invalid template",
			args{
				invalidTmplData(1),
				"../testdata/invalid.tmpl",
			},
			nil,
			"",
			true,
		},
		{
			"Default values",
			args{
				nil,
				"../testdata/tmpl_test.yml",
			},
			nil,
			yamlTemplateOutDefault,
			false,
		},
		{
			"Env values",
			args{
				nil,
				"../testdata/tmpl_test.yml",
			},
			map[string]string{
				"TEST_PGDBNAME": "foo",
				"TEST_PGUSER":   "bar",
			},
			yamlTemplateOutFromEnv,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				if err := os.Setenv(k, v); err != nil {
					t.Fatal(err)
				}
			}

			buf, err := yamlTemplate(tt.args.data, tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("yamlTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != "" {
				if got := buf.String(); got != tt.want {
					t.Errorf("yamlTemplate() =\n%v\nwant\n%v", got, tt.want)
				}
			}
		})
	}
}
