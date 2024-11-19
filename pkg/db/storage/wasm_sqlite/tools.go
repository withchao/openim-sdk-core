package wasm_sqlite

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"syscall/js"
)

const (
	funcOpen  = "sqliteOpen"
	funcClose = "sqliteClose"
	funcExec  = "sqliteExec"
	funcQuery = "sqliteQuery"
)

func waitAsyncFunc(ctx context.Context, result js.Value) (js.Value, error) {
	if !(result.Type() == js.TypeObject && result.InstanceOf(js.Global().Get("Promise"))) {
		return result, nil
	}
	var (
		resArgs []js.Value
		ok      bool
		done    = make(chan struct{})
	)
	success := js.FuncOf(func(this js.Value, args []js.Value) any {
		resArgs = args
		ok = true
		close(done)
		return nil
	})
	defer success.Release()
	failed := js.FuncOf(func(this js.Value, args []js.Value) any {
		resArgs = args
		close(done)
		return nil
	})
	defer failed.Release()
	result.Call("then", success).Call("catch", failed)
	select {
	case <-ctx.Done():
		return js.Null(), context.Cause(ctx)
	case <-done:
	}
	if !ok {
		return js.Null(), jsError(resArgs[0])
	}
	return resArgs[0], nil
}

func jsError(value js.Value) error {
	msg := value.Call("toString").String()
	log.Println("====>", msg)
	return fmt.Errorf("wasm return %w", errors.New(msg))
}

func call(ctx context.Context, name string, args ...any) (js.Value, error) {
	return waitAsyncFunc(ctx, js.Global().Call(name, args...))
}

func query(ctx context.Context, id int, name string, query string, args []driver.Value, resp any) error {
	if args == nil {
		args = []driver.Value{}
	}
	argsData, err := json.Marshal(args)
	if err != nil {
		return err
	}
	data, err := call(ctx, name, js.ValueOf(id), js.ValueOf(query), js.ValueOf(string(argsData)))
	if err != nil {
		return err
	}
	text := data.String()
	log.Println("#############query", "ctx", ctx.Value("ctx_test"), "sql", query, "args", string(argsData), "resp", text)
	if resp != nil {
		decoder := json.NewDecoder(bytes.NewReader([]byte(text)))
		decoder.UseNumber()
		if err := decoder.Decode(resp); err != nil {
			return fmt.Errorf("response %s unmarshal %w", text, err)
		}
	}
	return nil
}
