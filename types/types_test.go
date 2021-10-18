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

	"github.com/jackc/pgtype"
)

func TestNewNullGenerator(t *testing.T) {
	type args struct {
		seed            int64
		nullProbability int
	}
	tests := []struct {
		name string
		args args
		want *Probability
	}{
		{
			"Zero probability",
			args{
				1,
				0,
			},
			nil,
		},
		{
			"Negative probability",
			args{
				1,
				-2,
			},
			nil,
		},
		{
			"Positive probabilty",
			args{
				1,
				50,
			},
			NewProbability(1, 50),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNullGenerator(tt.args.seed, tt.args.nullProbability); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNullGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testType struct {
	pgtype.Int4
	nextVal int32
}

func (g *testType) NextValue() {
	g.Int4.Int = g.nextVal
	g.Int4.Status = pgtype.Present
}

func Test_valueGenerator_Get(t *testing.T) {
	type fields struct {
		ValueGenerator ValueGenerator
		nulls          *Probability
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			"null",
			fields{
				ValueGenerator: &testType{
					Int4: pgtype.Int4{
						Int:    3,
						Status: pgtype.Present,
					},
					nextVal: 22,
				},
				nulls: NewNullGenerator(1, 100),
			},
			nil,
		},
		{
			"nil nulls",
			fields{

				ValueGenerator: &testType{
					Int4: pgtype.Int4{
						Int:    3,
						Status: pgtype.Present,
					},
					nextVal: 22,
				},
				nulls: nil,
			},
			int32(22),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := valueGenerator{
				ValueGenerator: tt.fields.ValueGenerator,
				nulls:          tt.fields.nulls,
			}
			if got := v.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("valueGenerator.Get() = %T(%v), want %T(%v)", got, got, tt.want, tt.want)
			}
		})
	}
}

func Test_valueGenerator_AssignTo(t *testing.T) {
	type fields struct {
		ValueGenerator ValueGenerator
		nulls          *Probability
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			"null",
			fields{
				ValueGenerator: &testType{
					Int4: pgtype.Int4{
						Int:    3,
						Status: pgtype.Present,
					},
					nextVal: 22,
				},
				nulls: NewNullGenerator(1, 100),
			},
			0,
		},
		{
			"nil nulls",
			fields{

				ValueGenerator: &testType{
					Int4: pgtype.Int4{
						Int:    3,
						Status: pgtype.Present,
					},
					nextVal: 22,
				},
				nulls: nil,
			},
			22,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := valueGenerator{
				ValueGenerator: tt.fields.ValueGenerator,
				nulls:          tt.fields.nulls,
			}
			var got int

			v.AssignTo(&got)
			if got != tt.want {
				t.Errorf("valueGenerator.Get() = %T(%v), want %T(%v)", got, got, tt.want, tt.want)
			}
		})
	}
}
