package wasm_pglite

import (
	"context"
	"database/sql/driver"
	"log"
	"syscall/js"
)

type DriverContext struct{}

func (d DriverContext) open(name string) (int, error) {
	resp, err := call(context.Background(), "open", js.ValueOf(name))
	if err != nil {
		return 0, err
	}
	return resp.Int(), nil
}

func (d DriverContext) Open(name string) (driver.Conn, error) {
	log.Println("Driver.Open", name)
	id, err := d.open(name)
	if err != nil {
		return nil, err
	}
	return &Conn{id: id}, nil
}

func (d DriverContext) OpenConnector(name string) (driver.Connector, error) {
	//log.Println("DriverContext.OpenConnector", name)
	//id, err := d.open(name)
	//if err != nil {
	//	return nil, err
	//}
	return &Connector{name: name}, nil
}

type Connector struct {
	name string
}

func (c Connector) Connect(ctx context.Context) (driver.Conn, error) {
	return DriverContext{}.Open(c.name)
}

func (c Connector) Driver() driver.Driver {
	return DriverContext{}
}

type Conn struct {
	id int
}

func (c Conn) Prepare(query string) (driver.Stmt, error) {
	log.Println("Conn.Prepare", c.id, query)
	return &CustomStmt{id: c.id, query: query}, nil
}

func (c Conn) Close() error {
	log.Println("Conn.Close", c.id)
	_, err := call(context.Background(), "close", js.ValueOf(c.id))
	return err
}

func (c Conn) Begin() (driver.Tx, error) {
	log.Println("Conn.Begin", c.id)
	if err := query(context.Background(), c.id, "BEGIN", nil, nil); err != nil {
		return nil, err
	}
	return &Tx{id: c.id}, nil
}
