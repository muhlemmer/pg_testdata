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

// Package parse loads yaml from a configuration file, parses its arguments
// and initializes queries with type specific generator arguments.
package parse

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config structure root, meant to be marshalled / unmarshalled with yaml.
type Config struct {
	DSN    string // Data Source Name, aka connection string.
	Tables []*Table
}

// Load a yaml config file.
func Load(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("config.Load: %w", err)
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)

	conf := new(Config)

	if err = dec.Decode(conf); err != nil {
		return nil, fmt.Errorf("config.Load: %w", err)
	}

	conf.setColTableNames()

	return conf, nil
}

func (c *Config) setColTableNames() {
	for _, table := range c.Tables {
		table.setColTableNames()
	}
}
