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
	"text/template"
)

const insertQuery = "insert into {{ .Table }} ({{ .Columns }}) values ({{ .Positions }});"

var insertTmpl = template.Must(template.New("insertQuery").Parse(insertQuery))

type commaList []string

func (l commaList) String() string {
	return strings.Join([]string(l), ", ")
}

type insertData struct {
	Table     string
	Columns   commaList
	Positions commaList
}

// Table definition
type Table struct {
	Name    string // Name of the Table
	Amount  int    // Amount of Rows to generate and insert
	Columns []*column
}

func (table *Table) insert(tmpl *template.Template) (string, []interface{}) {
	data := insertData{
		Table:     table.Name,
		Columns:   make(commaList, len(table.Columns)),
		Positions: make(commaList, len(table.Columns)),
	}

	args := make([]interface{}, len(table.Columns))
	var err error

	for i, col := range table.Columns {
		data.Columns[i] = col.Name
		data.Positions[i] = fmt.Sprintf("$%d", i+1)
		args[i] = col.valueGenerator()
	}

	var buf strings.Builder
	if err = tmpl.Execute(&buf, &data); err != nil {
		panic(err)
	}

	return buf.String(), args
}

// InsertQuery with args for this table.
// The returned stmt can be used as prepared statement.
// The returned args can be reused for each iteration of the prepared statement.
// On each access, the args generate a new value, corresponding to the Generator
// options passed for each column / type.
func (table *Table) InsertQuery() (stmt string, args []interface{}, err error) {
	defer func() {
		if err, _ = recover().(error); err != nil {
			err = fmt.Errorf("parse.InsertQuery: %w in table %s", err, table.Name)
		}
	}()

	stmt, args = table.insert(insertTmpl)
	return
}
