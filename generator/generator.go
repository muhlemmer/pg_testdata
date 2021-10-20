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

// Package generator provides pgtype specific value generators.
// Generation is driven by the pseudo-random number generator from `math/rand`.
// This allows for deterministic value generation, using a seed for each constructor.
// Note that this determinism is also affected by other parameters,
// such as a minimum or maximum value.
//
// Each constructor also takes an argument for percentage of probability
// for a SQL null with each newly generated value.
// If the probability is 0 or lower, random null generation is disabled.
// Use this mode for columns with constraint NOT NULL.
// If the probability is 100 or higher, only nulls are generated.
package generator

import (
	"github.com/jackc/pgtype"
)

// newNull returns a Probability generator if nullProbability > 0, nil otherwise.
func newNull(seed int64, nullProbability float32) *probability {
	if nullProbability > 0 {
		return newProbability(seed, nullProbability)
	}
	return nil
}

// Value generates a pgtype value on each read access.
type Value interface {
	pgtype.ValueTranscoder
	// NextValue populates the Value with a newly generated value.
	NextValue()
}

type value struct {
	Value
	nulls *probability
}

func (v *value) nextStatusValue() {
	if v.nulls != nil && v.nulls.get() {
		v.Set(nil)
		return
	}
	v.NextValue()
}

func (v *value) AssignTo(dst interface{}) error {
	v.nextStatusValue()
	return v.Value.AssignTo(dst)
}

func (v value) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	v.nextStatusValue()
	return v.Value.EncodeBinary(ci, buf)
}

func (v value) EncodeText(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	v.nextStatusValue()
	return v.Value.EncodeText(ci, buf)
}

func (v value) Get() interface{} {
	v.nextStatusValue()
	return v.Value.Get()
}
