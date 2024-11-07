package wasm_pglite

import (
	"database/sql/driver"
	"fmt"
	"io"
	"time"
)

type rawRows struct {
	AffectedRows int              `json:"affectedRows"`
	Fields       []Field          `json:"fields"`
	Rows         []map[string]any `json:"rows"`
	columns      []string
	fieldType    map[string]int
	index        int
}

func (r *rawRows) init() {
	r.columns = make([]string, len(r.Fields))
	r.fieldType = make(map[string]int)
	for i, field := range r.Fields {
		r.columns[i] = field.Name
		r.fieldType[field.Name] = field.DataTypeID
	}
}

func (r *rawRows) Columns() []string {
	return r.columns
}

func (r *rawRows) Close() error {
	return nil
}

func (r *rawRows) toData(row map[string]any, field string) (driver.Value, error) {
	val, ok := row[field]
	if !ok {
		return nil, fmt.Errorf("not found field %s", field)
	}
	fieldType, ok := r.fieldType[field]
	if !ok {
		return nil, fmt.Errorf("undefined field %s", field)
	}
	if val == nil {
		return nil, nil
	}
	switch fieldType {
	case FieldTypeInt2:
		fallthrough
	case FieldTypeInt4:
		fallthrough
	case FieldTypeInt8:
		switch num := val.(type) {
		case float32:
			return int64(num), nil
		case float64:
			return int64(num), nil
		}
	case FieldTypeTimestamptz:
		return time.Parse(time.RFC3339, val.(string))
	}
	return val, nil
}

func (r *rawRows) Next(dest []driver.Value) error {
	index := r.index
	if len(r.Rows) <= index {
		return io.EOF
	}
	r.index++
	row := r.Rows[index]
	for i, column := range r.columns {
		//log.Println("rawRows.for====>", i, "/", column, ":", row[column], "->", printType(dest[i]))
		val, err := r.toData(row, column)
		if err != nil {
			return err
		}
		dest[i] = val
	}
	return nil
}

type Field struct {
	Name       string `json:"name"`
	DataTypeID int    `json:"dataTypeID"`
}

type baseExecResult struct {
	ID       int64 `json:"lastInsertId"`
	Affected int64 `json:"rowsAffected"`
}

func (r baseExecResult) RowsAffected() (int64, error) {
	return r.ID, nil
}

func (r baseExecResult) LastInsertId() (int64, error) {
	return r.Affected, nil
}
