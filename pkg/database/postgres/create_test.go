package postgres

import (
	"fmt"
	"testing"

	schemasv1alpha4 "github.com/schemahero/schemahero/pkg/apis/schemas/v1alpha4"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateTableStatement(t *testing.T) {
	tests := []struct {
		name              string
		tableSchema       *schemasv1alpha4.SQLTableSchema
		tableName         string
		expectedStatement string
	}{
		{
			name: "simple",
			tableSchema: &schemasv1alpha4.SQLTableSchema{
				PrimaryKey: []string{
					"id",
				},
				Columns: []*schemasv1alpha4.SQLTableColumn{
					&schemasv1alpha4.SQLTableColumn{
						Name: "id",
						Type: "integer",
					},
				},
			},
			tableName:         "simple",
			expectedStatement: `create table "simple" ("id" integer, primary key ("id"))`,
		},
		{
			name: "composite primary key",
			tableSchema: &schemasv1alpha4.SQLTableSchema{
				PrimaryKey: []string{
					"one",
					"two",
				},
				Columns: []*schemasv1alpha4.SQLTableColumn{
					&schemasv1alpha4.SQLTableColumn{
						Name: "one",
						Type: "integer",
					},
					&schemasv1alpha4.SQLTableColumn{
						Name: "two",
						Type: "integer",
					},
					&schemasv1alpha4.SQLTableColumn{
						Name: "three",
						Type: "varchar(255)",
					},
				},
			},
			tableName:         "composite_primary_key",
			expectedStatement: `create table "composite_primary_key" ("one" integer, "two" integer, "three" character varying (255), primary key ("one", "two"))`,
		},
		{
			name: "composite unique index",
			tableSchema: &schemasv1alpha4.SQLTableSchema{
				PrimaryKey: []string{
					"one",
				},
				Indexes: []*schemasv1alpha4.SQLTableIndex{
					&schemasv1alpha4.SQLTableIndex{
						Columns:  []string{"two", "three"},
						IsUnique: true,
					},
				},
				Columns: []*schemasv1alpha4.SQLTableColumn{
					&schemasv1alpha4.SQLTableColumn{
						Name: "one",
						Type: "integer",
					},
					&schemasv1alpha4.SQLTableColumn{
						Name: "two",
						Type: "integer",
					},
					&schemasv1alpha4.SQLTableColumn{
						Name: "three",
						Type: "varchar(255)",
					},
				},
			},
			tableName:         "composite_unique_index",
			expectedStatement: `create table "composite_unique_index" ("one" integer, "two" integer, "three" character varying (255), primary key ("one"), constraint "idx_composite_unique_index_two_three" unique ("two", "three"))`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := require.New(t)

			createTableStatement, err := CreateTableStatement(test.tableName, test.tableSchema)
			req.NoError(err)

			assert.Equal(t, test.expectedStatement, createTableStatement)
		})
	}
}

func Test_createExtensionStatement(t *testing.T) {
	schema := "public"
	version := "2.0.3"
	tests := []struct {
		name              string
		extention         *schemasv1alpha4.PostgreSQLExtension
		expectedStatement string
	}{
		{
			name: "simple",
			extention: &schemasv1alpha4.PostgreSQLExtension{
				Name:    "citext",
			},
			expectedStatement: `create extension "citext" if not exists`,
		},
		{
			name: "with_schema",
			extention: &schemasv1alpha4.PostgreSQLExtension{
				Name:    "hstore",
				Schema:  &schema,
			},
			expectedStatement: fmt.Sprintf(`create extension "hstore" if not exists schema "%s"`, schema),
		},
		{
			name: "with_schema_and_version",
			extention: &schemasv1alpha4.PostgreSQLExtension{
				Name:    "postgres_fdw",
				Schema:  &schema,
				Version: &version,
			},
			expectedStatement: fmt.Sprintf(`create extension "postgres_fdw" if not exists version "%s"`, version),
		},
		{
			name: "force_creation",
			extention: &schemasv1alpha4.PostgreSQLExtension{
				Name:    "postgis",
				Force:   true,
			},
			expectedStatement: `create extension "postgis"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := require.New(t)

			createTableStatement, err := createExtensionStatement(test.extention)
			req.NoError(err)

			assert.Equal(t, test.expectedStatement, createTableStatement)
		})
	}
}
