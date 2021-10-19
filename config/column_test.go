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

package config

import (
	"reflect"
	"testing"

	"github.com/muhlemmer/pg_testdata/types"
)

func TestMissingArgsError_Error(t *testing.T) {
	e := &MissingArgsError{
		Keys: []ArgName{MinArg, MaxArg},
		Type: Int4Type,
	}

	const want = "missing [min max] arguments for int4"

	if got := e.Error(); got != want {
		t.Errorf("MissingArgsError.Error() = %v, want %v", got, want)
	}

}

func TestColumn_requiredGenOpts(t *testing.T) {
	type args struct {
		tp   TypeName
		keys []ArgName
	}
	tests := []struct {
		name      string
		Generator map[ArgName]interface{}
		args      args
		wantErr   bool
	}{
		{
			"Panic",
			map[ArgName]interface{}{MinArg: 2},
			args{Int4Type, []ArgName{MinArg, MaxArg}},
			true,
		},
		{
			"Ok",
			map[ArgName]interface{}{MinArg: 2, MaxArg: 4},
			args{Int4Type, []ArgName{MinArg, MaxArg}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := func() (err error) {
				defer func() { err, _ = recover().(error) }()

				c := &Column{
					Generator: tt.Generator,
				}
				c.requiredGenOpts(tt.args.tp, tt.args.keys...)

				return nil
			}()

			if (err != nil) != tt.wantErr {
				t.Errorf("Column.requiredGenOpts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestColumn_boolType(t *testing.T) {
	type fields struct {
		Seed            int64
		NullProbability int
		Generator       map[ArgName]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		want    types.ValueGenerator
		wantErr bool
	}{
		{
			"Missing arg",
			fields{
				Seed:            1,
				NullProbability: 2,
				Generator:       nil,
			},
			nil,
			true,
		},
		{
			"Wrong probability type",
			fields{
				Seed:            1,
				NullProbability: 2,
				Generator:       map[ArgName]interface{}{ProbabilityArg: "foo"},
			},
			nil,
			true,
		},
		{
			"OK",
			fields{
				Seed:            1,
				NullProbability: 2,
				Generator:       map[ArgName]interface{}{ProbabilityArg: 50},
			},
			types.NewBool(1, 2, 50),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Column{
				Seed:            tt.fields.Seed,
				NullProbability: tt.fields.NullProbability,
				Generator:       tt.fields.Generator,
			}

			err := func() (err error) {
				defer func() { err, _ = recover().(error) }()
				if got := c.boolType(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Column.boolType() = %v, want %v", got, tt.want)
				}
				return nil
			}()

			if (err != nil) != tt.wantErr {
				t.Errorf("Column.boolType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

const unsupportedType TypeName = "unsupported"

func TestColumn_ValueGenerator(t *testing.T) {
	column := Column{
		Name:            "test_column",
		Seed:            1,
		NullProbability: 2,
		table: &Table{
			Name: "test_table",
		},
	}

	type fields struct {
		Type      TypeName
		Generator map[ArgName]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		wantVg  types.ValueGenerator
		wantErr bool
	}{
		{
			"unsupported type",
			fields{
				Type:      unsupportedType,
				Generator: nil,
			},
			nil,
			true,
		},
		{
			"bool type",
			fields{
				Type:      BoolType,
				Generator: map[ArgName]interface{}{ProbabilityArg: 50},
			},
			types.NewBool(1, 2, 50),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Column{
				Name:            column.Name,
				Seed:            column.Seed,
				NullProbability: column.NullProbability,
				Type:            tt.fields.Type,
				Generator:       tt.fields.Generator,
				table:           column.table,
			}
			gotVg, err := c.ValueGenerator()
			if (err != nil) != tt.wantErr {
				t.Errorf("Column.ValueGenerator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVg, tt.wantVg) {
				t.Errorf("Column.ValueGenerator() = %v, want %v", gotVg, tt.wantVg)
			}
		})
	}
}
