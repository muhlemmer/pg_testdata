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
	"fmt"
	"reflect"
	"testing"

	"github.com/muhlemmer/pg_testdata/generator"
)

func Test_columnError_Error(t *testing.T) {
	e := &columnError{
		fmt.Errorf("foobar"),
		"column",
	}

	const want = "foobar in column \"column\""

	if got := e.Error(); got != want {
		t.Errorf("columnError.Error() = %v, want %v", got, want)
	}

}

func Test_column_requiredGenOpts(t *testing.T) {
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
				t.Errorf("column.requiredGenOpts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_assertFloat32(t *testing.T) {
	c := &Column{
		Name: "test",
	}

	tests := []struct {
		name    string
		v       interface{}
		want    float32
		wantErr bool
	}{
		{
			"float32",
			float32(1),
			1,
			false,
		},
		{
			"float64",
			float64(1),
			1,
			false,
		},
		{
			"int",
			1,
			1,
			false,
		},
		{
			"string",
			"foo",
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := func() (err error) {
				defer func() { err, _ = recover().(error) }()

				if got := c.assertFloat32(tt.v); got != tt.want {
					t.Errorf("column.assertFloat32() = %v, want %v", got, tt.want)
				}

				return
			}()

			if (err != nil) != tt.wantErr {
				t.Errorf("column.boolType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_column_boolType(t *testing.T) {
	type fields struct {
		Seed            int64
		NullProbability float32
		Generator       map[ArgName]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		want    generator.Value
		wantErr bool
	}{
		{
			"Missing arg",
			fields{
				Seed:            1,
				NullProbability: 2.0,
				Generator:       nil,
			},
			nil,
			true,
		},
		{
			"Wrong probability type",
			fields{
				Seed:            1,
				NullProbability: 2.0,
				Generator:       map[ArgName]interface{}{ProbabilityArg: "foo"},
			},
			nil,
			true,
		},
		{
			"OK",
			fields{
				Seed:            1,
				NullProbability: 2.0,
				Generator:       map[ArgName]interface{}{ProbabilityArg: float32(50.0)},
			},
			generator.NewBool(1, 2, 50),
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
				t.Errorf("column.boolType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

const unsupportedType TypeName = "unsupported"

func Test_column_valueGenerator(t *testing.T) {
	col := Column{
		Name:            "test_column",
		Seed:            1,
		NullProbability: 2,
	}

	type fields struct {
		Type      TypeName
		Generator map[ArgName]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		wantVg  generator.Value
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
				Generator: map[ArgName]interface{}{ProbabilityArg: float32(50)},
			},
			generator.NewBool(1, 2, 50),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := func() (err error) {
				defer func() { err, _ = recover().(error) }()

				c := &Column{
					Name:            col.Name,
					Seed:            col.Seed,
					NullProbability: col.NullProbability,
					Type:            tt.fields.Type,
					Generator:       tt.fields.Generator,
				}
				gotVg := c.valueGenerator()
				if !reflect.DeepEqual(gotVg, tt.wantVg) {
					t.Errorf("column.valueGenerator() = %v, want %v", gotVg, tt.wantVg)
				}

				return
			}()

			if (err != nil) != tt.wantErr {
				t.Errorf("column.valueGenerator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
