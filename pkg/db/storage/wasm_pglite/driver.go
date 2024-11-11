package wasm_pglite

import (
	"context"
	"database/sql/driver"
	"log"
	"strings"
	"syscall/js"
)

type DriverContext struct{}

func (d DriverContext) open(name string) (int, error) {
	resp, err := call(context.Background(), funcOpen, js.ValueOf(name))
	if err != nil {
		return 0, err
	}
	log.Println("#############open", resp.Type().String())
	return resp.Int(), nil
}

func (d DriverContext) Open(name string) (driver.Conn, error) {
	id, err := d.open(name)
	if err != nil {
		return nil, err
	}
	return &Conn{id: id}, nil
}

func (d DriverContext) OpenConnector(name string) (driver.Connector, error) {
	return &Connector{name: name}, nil
}

type Connector struct {
	name string
}

func (c Connector) Connect(ctx context.Context) (driver.Conn, error) {
	return (DriverContext{}).Open(c.name)
}

func (c Connector) Driver() driver.Driver {
	return DriverContext{}
}

type Conn struct {
	id int
}

func (c Conn) Prepare(query string) (driver.Stmt, error) {
	stmt := &Stmt{id: c.id, query: query}
	if query != autoMigrateSQL {
		stmt.query = strings.ReplaceAll(stmt.query, "`", "")
		return stmt, nil
	}
	return pgliteAutoMigrateCustomStmt{
		Stmt: stmt,
	}, nil
}

func (c Conn) Close() error {
	_, err := call(context.Background(), funcClose, js.ValueOf(c.id))
	return err
}

func (c Conn) Begin() (driver.Tx, error) {
	if err := query(context.Background(), c.id, "BEGIN", nil, nil); err != nil {
		return nil, err
	}
	return &Tx{id: c.id}, nil
}
