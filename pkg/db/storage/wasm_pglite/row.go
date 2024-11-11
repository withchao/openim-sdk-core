package wasm_pglite

import (
	"database/sql/driver"
	"fmt"
	"io"
	"time"
)

const (
	FieldTypeInt2 = 21 // smallint
	FieldTypeInt4 = 23 // integer
	FieldTypeInt8 = 20 // bigint

	FieldTypeFloat4  = 700  // real
	FieldTypeFloat8  = 701  // double precision
	FieldTypeNumeric = 1700 // numeric

	FieldTypeChar    = 18   // char
	FieldTypeVarchar = 1043 // varchar
	FieldTypeText    = 25   // text
	FieldTypeName    = 19   // name

	FieldTypeBool = 16 // bool

	FieldTypeDate        = 1082 // date
	FieldTypeTime        = 1083 // time
	FieldTypeTimestamp   = 1114 // timestamp
	FieldTypeTimestamptz = 1184 // timestamp with time zone
	FieldTypeInterval    = 1186 // interval

	FieldTypePoint   = 600 // point
	FieldTypeLine    = 628 // line
	FieldTypeLseg    = 601 // lseg
	FieldTypeBox     = 603 // box
	FieldTypePath    = 602 // path
	FieldTypePolygon = 604 // polygon
	FieldTypeCircle  = 718 // circle

	FieldTypeCIDR    = 650 // cidr
	FieldTypeINET    = 869 // inet
	FieldTypeMacAddr = 829 // macaddr

	FieldTypeJSON  = 114  // json
	FieldTypeJSONB = 3802 // jsonb

	FieldTypeUUID = 2950 // uuid

	FieldTypeInt4Array = 1007 // int4[]
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
	Affected int64 `json:"affectedRows"`
}

func (r baseExecResult) RowsAffected() (int64, error) {
	return r.Affected, nil
}

func (r baseExecResult) LastInsertId() (int64, error) {
	return r.ID, nil
}
