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
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func lookupEnv(key, defValue string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return defValue
}

var funcMap = template.FuncMap{"env": lookupEnv}

func yamlTemplate(data interface{}, filename string) (*bytes.Buffer, error) {
	tmpl := template.New(filepath.Base(filename))
	tmpl.Funcs(funcMap)

	var err error

	tmpl, err = tmpl.ParseFiles(filename)
	if err != nil {
		return nil, fmt.Errorf("parse.yamlTemplate: %w", err)
	}

	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, data); err != nil {
		return nil, fmt.Errorf("parse.yamlTemplate: %w", err)
	}

	return buf, nil
}
