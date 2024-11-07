package db

import (
	"context"
)

func (d *DataBase) GetExistTables(ctx context.Context) ([]string, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	tables, err := d.conn.WithContext(ctx).Migrator().GetTables()
	if err != nil {
		return nil, err
	}
	return tables, nil
}
