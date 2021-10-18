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
	"github.com/jackc/pgtype"
)

// NewNullGenerator returns a Probability if nullProbability > 0, nil otherwise.
func NewNullGenerator(seed int64, nullProbability int) *Probability {
	if nullProbability > 0 {
		return NewProbability(seed, nullProbability)
	}
	return nil
}

type ValueGenerator interface {
	pgtype.ValueTranscoder
	NextValue()
}

type valueGenerator struct {
	ValueGenerator
	nulls *Probability
}

func (v *valueGenerator) nextStatusValue() {
	if v.nulls != nil && v.nulls.Get() {
		v.Set(nil)
		return
	}
	v.NextValue()
}

func (v *valueGenerator) AssignTo(dst interface{}) error {
	v.nextStatusValue()
	return v.ValueGenerator.AssignTo(dst)
}

func (v valueGenerator) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	v.nextStatusValue()
	return v.ValueGenerator.EncodeBinary(ci, buf)
}

func (v valueGenerator) EncodeText(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	v.nextStatusValue()
	return v.ValueGenerator.EncodeText(ci, buf)
}

func (v valueGenerator) Get() interface{} {
	v.nextStatusValue()
	return v.ValueGenerator.Get()
}
