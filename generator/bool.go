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

package generator

import (
	"github.com/jackc/pgtype"
)

type boolType struct {
	pgtype.Bool
	generator *probability
}

func (b *boolType) NextValue() {
	b.Bool.Bool = b.generator.get()
	b.Bool.Status = pgtype.Present
}

// NewBool returns a boolean value generator.
//
// Probability is a percentage of chance `true` values are generated on each read.
// If probability is 0 or lower, only `false` values are generated.
// If probability is 100 or highter, only `true` values are generated.
func NewBool(seed int64, nullProbabilty, probabilty float32) Value {
	return &value{
		Value: &boolType{
			generator: newProbability(seed, probabilty),
		},
		nulls: newNull(seed, nullProbabilty),
	}
}
