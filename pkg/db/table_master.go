package db

import (
	"context"
	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) GetExistTables(ctx context.Context) ([]string, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	tables, err := d.conn.WithContext(ctx).Migrator().GetTables()
	if err != nil {
		return nil, errs.WrapMsg(err, "GetTables failed")
	}
	return tables, nil
}
