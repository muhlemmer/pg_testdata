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
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

var testConf = Config{
	DSN: "dbname=testdata user=testdata host=db port=5432 connect_timeout=10",
	Tables: []*Table{
		{
			Name:   "all_supported",
			Amount: 1000,
			MaxDuration: TableDurations{
				Table: time.Minute,
				Exec:  time.Second,
			},
			Columns: []*Column{
				{
					Name:            "bool_col_n",
					Seed:            2,
					NullProbability: 10.0,
					Type:            "bool",
					Generator: map[ArgName]interface{}{
						ProbabilityArg: 70.1,
					},
				},
				{
					Name:            "bool_col_nn",
					Seed:            2,
					NullProbability: 0.0,
					Type:            "bool",
					Generator: map[ArgName]interface{}{
						ProbabilityArg: 70.1,
					},
				},
			},
		},
	},
}

func writeExample(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("writeExample: %w", err)
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	defer enc.Close()

	if err = enc.Encode(&testConf); err != nil {
		return fmt.Errorf("writeExample: %w", err)
	}

	return nil
}
