package wasm_sqlite

import (
	"database/sql/driver"
	"encoding/json"
	"io"
	"strings"
)

type rawRows struct {
	Columns_    []string         `json:"columns"`
	Values      [][]driver.Value `json:"values"`
	index       int
	columnIndex map[string]int
}

func (r *rawRows) init() {
	r.columnIndex = make(map[string]int)
	for i, name := range r.Columns_ {
		r.columnIndex[name] = i
	}
}

func (r *rawRows) Columns() []string {
	return r.Columns_
}

func (r *rawRows) Close() error {
	return nil
}

func (r *rawRows) Next(dest []driver.Value) error {
	index := r.index
	if len(r.Values) <= index {
		return io.EOF
	}
	r.index++
	row := r.Values[index]
	for i := range dest {
		elem := row[i]
		if elem == nil {
			continue
		}
		if num, ok := elem.(json.Number); ok {
			var err error
			if strings.IndexByte(num.String(), '.') >= 0 {
				elem, err = num.Float64()
			} else {
				elem, err = num.Int64()
			}
			if err != nil {
				return err
			}
		}
		dest[i] = elem
	}
	return nil
}

type baseExecResult struct {
	Affected int64
}

func (r baseExecResult) RowsAffected() (int64, error) {
	return r.Affected, nil
}

func (r baseExecResult) LastInsertId() (int64, error) {
	return 0, nil
}
