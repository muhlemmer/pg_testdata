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
	"log"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	if err := writeExample("../testdata/all_supported.yml"); err != nil {
		log.Fatalln(err)
	}

	os.Exit(m.Run())
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     *Config
		wantErr  bool
	}{
		{
			"File not found",
			"/foo/bar/does/not/exist",
			nil,
			true,
		},
		{
			"Invalid file",
			"../testdata/invalid.yml",
			nil,
			true,
		},
		{
			"Example config",
			"../testdata/all_supported.yml",
			&testConf,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}
