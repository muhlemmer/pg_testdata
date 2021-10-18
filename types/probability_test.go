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

func allTrueBool(size int) []bool {
	bs := make([]bool, size)
	for i := range bs {
		bs[i] = true
	}

	return bs
}

func TestProbability_Get(t *testing.T) {
	tests := []struct {
		name string
		p    *Probability
		want []bool
	}{
		{
			"0 probability",
			NewProbability(1, 0),
			make([]bool, 100),
		},
		{
			"100 probability",
			NewProbability(1, 100),
			allTrueBool(100),
		},
		{
			"50 probability",
			NewProbability(1, 50),
			[]bool{false, false, true, false, false, true, true, true, false, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]bool, len(tt.want))

			for i := range got {
				got[i] = tt.p.Get()
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Probability.Get() =\n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}
