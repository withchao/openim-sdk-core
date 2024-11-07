package wasm_pglite

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
	var res baseExecResult
	if err := query(context.Background(), s.id, s.query, args, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	var res rawRows
	if err := query(context.Background(), s.id, s.query, args, &res); err != nil {
		return nil, err
	}
	res.init()
	return &res, nil
}

func (s *Stmt) NumInput() int {
	return -1
}
