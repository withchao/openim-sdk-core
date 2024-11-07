package wasm_pglite

import "context"

type Tx struct {
	id int
}

func (t Tx) Commit() error {
	return query(context.Background(), t.id, "COMMIT", nil, nil)
}

func (t Tx) Rollback() error {
	return query(context.Background(), t.id, "ROLLBACK", nil, nil)
}
