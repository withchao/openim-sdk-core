package wasm_sqlite

import (
	"context"
	"database/sql/driver"
)

type Stmt struct {
	id    int
	query string
}

func (s *Stmt) Close() error {
	return nil
}

func (s *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	var affected int64
	if err := query(context.Background(), s.id, funcExec, s.query, args, &affected); err != nil {
		return nil, err
	}
	return &baseExecResult{Affected: affected}, nil
}

func (s *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	var res []rawRows
	if err := query(context.Background(), s.id, funcQuery, s.query, args, &res); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &rawRows{}, nil
	}
	val := res[0]
	val.init()
	return &val, nil
}

func (s *Stmt) NumInput() int {
	return -1
}
