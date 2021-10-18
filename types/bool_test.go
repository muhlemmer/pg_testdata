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

package types

import (
	"reflect"
	"testing"
)

func Test_boolType(t *testing.T) {
	type args struct {
		seed            int64
		nullProbability int
		probability     int
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			"null",
			args{1, 100, 100},
			nil,
		},
		{
			"true",
			args{1, 0, 100},
			true,
		},
		{
			"false",
			args{1, 0, 0},
			false,
		},
		{
			"random",
			args{3, 0, 50},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewBool(tt.args.seed, tt.args.nullProbability, tt.args.probability)

			if got := v.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("valueGenerator.Get() = %T(%v), want %T(%v)", got, got, tt.want, tt.want)
			}
		})
	}
}
