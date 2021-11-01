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
	"strings"

	"github.com/muhlemmer/pg_testdata/generator"
)

type TypeName string

const (
	BoolType TypeName = "bool"
	Int4Type TypeName = "int4"
)

type ArgName string

const (
	MinArg         ArgName = "min"
	MaxArg         ArgName = "max"
	ProbabilityArg ArgName = "probability"
)

// Column information and parameters.
type Column struct {
	Name            string
	Seed            int64
	NullProbability float32
	Type            TypeName
	Generator       map[ArgName]interface{}
}

type columnError struct {
	err    error
	column string
}

func (e *columnError) Error() string {
	return fmt.Sprintf("%v in column %q", e.err, e.column)
}

func (c *Column) panic(err error) {
	panic(&columnError{
		err:    err,
		column: c.Name,
	})
}

// requiredGenOpts checks if the required "keys" are present in the
// Generator arguments map. If any keys are found missing,
// all missing keys are collected and passed to panic() in a MissingArgsError.
func (c *Column) requiredGenOpts(tp TypeName, keys ...ArgName) {
	var missing []string

	for _, k := range keys {
		if _, ok := c.Generator[k]; !ok {
			missing = append(missing, string(k))
		}
	}

	if len(missing) > 0 {
		c.panic(fmt.Errorf("missing arguments %q for type %q", strings.Join(missing, " ,"), tp))
	}
}

func (c *Column) assertFloat32(v interface{}) float32 {
	switch f := v.(type) {
	case float32:
		return f
	case float64:
		return float32(f)
	case int:
		return float32(f)
	default:
		c.panic(fmt.Errorf("bool \"probabilty\" incorrect type: %T, expected: float32", v))
		return 0
	}
}

func (c *Column) boolType() generator.Value {
	c.requiredGenOpts(BoolType, ProbabilityArg)

	return generator.NewBool(c.Seed, c.NullProbability, c.assertFloat32(c.Generator[ProbabilityArg]))
}

// valueGenerator panics in case of an invalid Type argument.
func (c *Column) valueGenerator() generator.Value {
	switch c.Type {
	case BoolType:
		return c.boolType()
	default:
		c.panic(fmt.Errorf("unsuported type %q", c.Type))
		return nil
	}
}
