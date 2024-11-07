package wasm_pglite

import (
	"database/sql/driver"
)

const autoMigrateSQL = `SELECT c.column_name, c.is_nullable = 'YES', c.udt_name, c.character_maximum_length, c.numeric_precision, c.numeric_precision_radix, c.numeric_scale, c.datetime_precision, 8 * typlen, c.column_default, pd.description, c.identity_increment FROM information_schema.columns AS c JOIN pg_type AS pgt ON c.udt_name = pgt.typname LEFT JOIN pg_catalog.pg_description as pd ON pd.objsubid = c.ordinal_position AND pd.objoid = (SELECT oid FROM pg_catalog.pg_class WHERE relname = c.table_name AND relnamespace = (SELECT oid FROM pg_catalog.pg_namespace WHERE nspname = c.table_schema)) where table_catalog = $1 AND table_schema = CURRENT_SCHEMA() AND table_name = $2`

type pgliteAutoMigrateCustomStmt struct {
	*Stmt
}

func (s pgliteAutoMigrateCustomStmt) Query(args []driver.Value) (driver.Rows, error) {
	res, err := s.Stmt.Query(args)
	if err != nil {
		return nil, err
	}
	raws, ok := res.(*rawRows)
	if !ok {
		return res, nil
	}
	const (
		column  = "?column?"
		column1 = "?column?(1)"
	)
	if _, ok := raws.fieldType[column]; !ok {
		return res, nil
	}
	if _, ok := raws.fieldType[column1]; ok {
		return res, nil
	}
	for i, field := range raws.Fields {
		if field.Name == column {
			raws.Fields[i].Name = column1
			break
		}
	}
	raws.Fields = append(raws.Fields, Field{
		Name:       column,
		DataTypeID: FieldTypeBool,
	})
	for i, row := range raws.Rows {
		raws.Rows[i][column1] = row[column]
		raws.Rows[i][column] = 1
	}
	return raws, nil
}
