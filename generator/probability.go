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
	"math/rand"
)

// maxProbability represents 100%.
const maxProbability = 100.0

// probability is a pseudo-random and deterministic generator of booleans.
// With a probability value of 100 or higher, only `true` will be generated
// and a probability value of 0 or lower will only generate `false`.
type probability struct {
	rand        *rand.Rand
	probability float32
}

// newProbability returns a Probability, with the random source initialized with seed.
// and the percentage of probability at which it will generate `true` values.
func newProbability(seed int64, prob float32) *probability {
	return &probability{
		rand: rand.New(
			rand.NewSource(seed),
		),
		probability: prob,
	}
}

// get the next random bool value.
func (p *probability) get() bool {
	return p.rand.Float32()*100.0 < p.probability
}
