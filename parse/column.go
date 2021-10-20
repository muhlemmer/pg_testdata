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

	"github.com/muhlemmer/pg_testdata/types"
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
	NullProbability int
	Type            TypeName
	Generator       map[ArgName]interface{}

	table *Table
}

// MissingArgsError is returned when required Generator arguments
// for the specified column type are missing
type MissingArgsError struct {
	Keys []ArgName
	Type TypeName
}

func (e *MissingArgsError) Error() string {
	return fmt.Sprintf("missing %s arguments for %s", e.Keys, e.Type)
}

// requiredGenOpts checks if the required "keys" are present in the
// Generator arguments map. If any keys are found missing,
// all missing keys are collected and passed to panic() in a MissingArgsError.
func (c *Column) requiredGenOpts(tp TypeName, keys ...ArgName) {
	var missing []ArgName

	for _, k := range keys {
		if _, ok := c.Generator[k]; !ok {
			missing = append(missing, k)
		}
	}

	if len(missing) > 0 {
		panic(&MissingArgsError{missing, tp})
	}
}

func (c *Column) boolType() types.ValueGenerator {
	c.requiredGenOpts(BoolType, ProbabilityArg)

	prob, ok := c.Generator[ProbabilityArg].(int)
	if !ok {
		panic(fmt.Errorf("bool \"probabilty\" incorrect type: %T, expected: int", c.Generator[ProbabilityArg]))
	}

	return types.NewBool(c.Seed, c.NullProbability, prob)
}

// valueGenerator panics in case of an invalid Type argument.
func (c *Column) valueGenerator() types.ValueGenerator {
	switch c.Type {
	case BoolType:
		return c.boolType()
	default:
		panic(fmt.Errorf("unsuported type %q", c.Type))
	}
}

// ValueGenerator constructs the generator for this column.
// An error is returned when the requested Type is not supported,
// any required Generator arguments missing or a wrong type of any of the Generator arguments.
func (c *Column) ValueGenerator() (vg types.ValueGenerator, err error) {
	defer func() {
		err, _ = recover().(error)
		if err != nil {
			err = fmt.Errorf("ValueGenerator: %w for column %s.%s", err, c.table.Name, c.Name)
		}
	}()

	return c.valueGenerator(), nil
}
