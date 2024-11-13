package wasm_sqlite

import "context"

type Tx struct {
	id int
}

func (t Tx) Commit() error {
	return query(context.Background(), t.id, funcExec, "COMMIT", nil, nil)
}

func (t Tx) Rollback() error {
	return query(context.Background(), t.id, funcExec, "ROLLBACK", nil, nil)
}
