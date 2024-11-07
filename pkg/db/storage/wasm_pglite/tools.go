package wasm_pglite

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"strconv"
	"syscall/js"
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

func jsError(value js.Value) (err error) {
	if value.Type() == js.TypeObject {
		code := value.Get("code")
		msg := value.Get("message")
		if code.Type() == js.TypeString && msg.Type() == js.TypeString {
			getStr := func(name string) string {
				val := value.Get(name)
				if val.Type() == js.TypeString {
					return val.String()
				} else {
					return ""
				}
			}
			getInt32 := func(name string) int32 {
				val := value.Get(name)
				switch val.Type() {
				case js.TypeNumber:
					return int32(val.Int())
				case js.TypeString:
					tmp, _ := strconv.Atoi(val.String())
					return int32(tmp)
				default:
					return 0
				}
			}
			return &pgconn.PgError{
				Severity:         getStr("severity"),
				Code:             code.String(),
				Message:          msg.String(),
				Detail:           getStr("detail"),
				Hint:             getStr("hint"),
				Position:         getInt32("position"),
				InternalPosition: getInt32("internalPosition"),
				InternalQuery:    getStr("internalQuery"),
				Where:            getStr("where"),
				SchemaName:       getStr("schema"),
				TableName:        getStr("table"),
				ColumnName:       getStr("column"),
				DataTypeName:     getStr("dataType"),
				ConstraintName:   getStr("constraint"),
				File:             getStr("file"),
				Line:             getInt32("line"),
				Routine:          getStr("heap_create_with_catalog"),
			}
		}
	}
	switch value.Type() {
	case js.TypeString:
		return errors.New(value.String())
	default:
		return fmt.Errorf("unknown js error %v", value)
	}
}

func call(ctx context.Context, name string, args ...any) (js.Value, error) {
	return waitAsyncFunc(ctx, js.Global().Get("gopglite").Call(name, args...))
}

func query(ctx context.Context, id int, query string, args []driver.Value, resp any) error {
	if args == nil {
		args = []driver.Value{}
	}
	argsData, err := json.Marshal(args)
	if err != nil {
		return err
	}
	data, err := call(ctx, "query", js.ValueOf(id), js.ValueOf(query), js.ValueOf(string(argsData)))
	if err != nil {
		return err
	}
	if resp != nil {
		text := data.String()
		if err := json.Unmarshal([]byte(text), resp); err != nil {
			return fmt.Errorf("response %s unmarshal %w", text, err)
		}
	}
	return nil
}
