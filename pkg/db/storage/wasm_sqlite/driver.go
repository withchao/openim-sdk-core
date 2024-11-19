package wasm_sqlite

import (
	"context"
	"database/sql/driver"
	"syscall/js"
)

type DriverContext struct {
	ctx context.Context
}

func (d DriverContext) open(name string) (int, error) {
	resp, err := call(d.ctx, funcOpen, js.ValueOf(name))
	if err != nil {
		return 0, err
	}
	return resp.Int(), nil
}

func (d DriverContext) Open(name string) (driver.Conn, error) {
	id, err := d.open(name)
	if err != nil {
		return nil, err
	}
	return &Conn{ctx: d.ctx, id: id}, nil
}

func (d DriverContext) OpenConnector(name string) (driver.Connector, error) {
	return &Connector{name: name}, nil
}

type Connector struct {
	name string
}

func (c Connector) Connect(ctx context.Context) (driver.Conn, error) {
	return (DriverContext{ctx: ctx}).Open(c.name)
}

func (c Connector) Driver() driver.Driver {
	return DriverContext{}
}

type Conn struct {
	ctx context.Context
	id  int
}

func (c Conn) Prepare(query string) (driver.Stmt, error) {
	return &Stmt{ctx: c.ctx, id: c.id, query: query}, nil
}

func (c Conn) Close() error {
	_, err := call(c.ctx, funcClose, js.ValueOf(c.id))
	return err
}

func (c Conn) Begin() (driver.Tx, error) {
	if err := query(c.ctx, c.id, funcExec, "BEGIN", nil, nil); err != nil {
		return nil, err
	}
	return &Tx{ctx: c.ctx, id: c.id}, nil
}
