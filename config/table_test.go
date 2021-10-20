package config

import (
	"reflect"
	"testing"
	"text/template"

	"github.com/muhlemmer/pg_testdata/types"
)

func Test_commaList_String(t *testing.T) {
	tests := []struct {
		name string
		l    commaList
		want string
	}{
		{
			"One field",
			commaList{"one"},
			"one",
		},
		{
			"Two fields",
			commaList{"one", "two"},
			"one, two",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("commaList.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Table_insert(t *testing.T) {
	errTmpl := template.Must(template.New("foo").Parse("insert {{ .Foo }}"))

	tests := []struct {
		name    string
		table   Table
		tmpl    *template.Template
		want    string
		want1   []interface{}
		wantErr bool
	}{
		{
			"Generator error",
			Table{
				Name:   "articles",
				Amount: 10,
				Columns: []*Column{
					{
						Name:            "published",
						Seed:            1,
						NullProbability: 0,
						Type:            BoolType,
						Generator:       nil,
					},
				},
			},
			insertTmpl,
			"",
			nil,
			true,
		},
		{
			"Template error",
			Table{
				Name:   "articles",
				Amount: 10,
				Columns: []*Column{
					{
						Name:            "published",
						Seed:            1,
						NullProbability: 0,
						Type:            BoolType,
						Generator: map[ArgName]interface{}{
							ProbabilityArg: 1,
						},
					},
				},
			},
			errTmpl,
			"",
			nil,
			true,
		},
		{
			"Success",
			Table{
				Name:   "articles",
				Amount: 10,
				Columns: []*Column{
					{
						Name:            "published",
						Seed:            1,
						NullProbability: 0,
						Type:            BoolType,
						Generator: map[ArgName]interface{}{
							ProbabilityArg: 1,
						},
					},
					{
						Name:            "special",
						Seed:            2,
						NullProbability: 50,
						Type:            BoolType,
						Generator: map[ArgName]interface{}{
							ProbabilityArg: 99,
						},
					},
				},
			},
			insertTmpl,
			"insert into articles (published, special) values ($1, $2);",
			[]interface{}{
				types.NewBool(1, 0, 1),
				types.NewBool(2, 50, 99),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.table.setColTableNames()

			got, got1, err := tt.table.insert(tt.tmpl)
			if (err != nil) != tt.wantErr {
				t.Errorf("Table.insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Table.insert() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Table.insert() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_Table_Insert(t *testing.T) {
	tests := []struct {
		name    string
		table   Table
		want    string
		want1   []interface{}
		wantErr bool
	}{
		{
			"Generator error",
			Table{
				Name:   "articles",
				Amount: 10,
				Columns: []*Column{
					{
						Name:            "published",
						Seed:            1,
						NullProbability: 0,
						Type:            BoolType,
						Generator:       nil,
					},
				},
			},
			"",
			nil,
			true,
		},
		{
			"Success",
			Table{
				Name:   "articles",
				Amount: 10,
				Columns: []*Column{
					{
						Name:            "published",
						Seed:            1,
						NullProbability: 0,
						Type:            BoolType,
						Generator: map[ArgName]interface{}{
							ProbabilityArg: 1,
						},
					},
					{
						Name:            "special",
						Seed:            2,
						NullProbability: 50,
						Type:            BoolType,
						Generator: map[ArgName]interface{}{
							ProbabilityArg: 99,
						},
					},
				},
			},
			"insert into articles (published, special) values ($1, $2);",
			[]interface{}{
				types.NewBool(1, 0, 1),
				types.NewBool(2, 50, 99),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.table.setColTableNames()

			got, got1, err := tt.table.Insert()
			if (err != nil) != tt.wantErr {
				t.Errorf("Table.insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Table.insert() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Table.insert() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
