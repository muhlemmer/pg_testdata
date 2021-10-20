package config

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
	Columns []*Column
}

func (table *Table) setColTableNames() {
	for i := range table.Columns {
		table.Columns[i].table = table
	}
}

func (table *Table) insert(tmpl *template.Template) (string, []interface{}, error) {
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

		args[i], err = col.ValueGenerator()
		if err != nil {
			return "", nil, err
		}
	}

	var buf strings.Builder
	if err = tmpl.Execute(&buf, &data); err != nil {
		return "", nil, err
	}

	return buf.String(), args, nil
}

func (table *Table) Insert() (stmt string, args []interface{}, err error) {
	if stmt, args, err = table.insert(insertTmpl); err != nil {
		err = fmt.Errorf("query.Insert: %w", err)
	}

	return
}
