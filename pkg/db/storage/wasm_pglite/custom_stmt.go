package wasm_pglite

import (
	"context"
	"database/sql/driver"
	"log"
)

const (
	// 整型
	FieldTypeInt2 = 21 // smallint
	FieldTypeInt4 = 23 // integer
	FieldTypeInt8 = 20 // bigint

	// 浮点型
	FieldTypeFloat4  = 700  // real
	FieldTypeFloat8  = 701  // double precision
	FieldTypeNumeric = 1700 // numeric

	// 字符型
	FieldTypeChar    = 18   // char
	FieldTypeVarchar = 1043 // varchar
	FieldTypeText    = 25   // text
	FieldTypeName    = 19   // name (数据库中的名称)

	// 布尔型
	FieldTypeBool = 16 // bool

	// 日期/时间型
	FieldTypeDate        = 1082 // date
	FieldTypeTime        = 1083 // time
	FieldTypeTimestamp   = 1114 // timestamp
	FieldTypeTimestamptz = 1184 // timestamp with time zone
	FieldTypeInterval    = 1186 // interval

	// 几何类型
	FieldTypePoint   = 600 // point
	FieldTypeLine    = 628 // line
	FieldTypeLseg    = 601 // lseg
	FieldTypeBox     = 603 // box
	FieldTypePath    = 602 // path
	FieldTypePolygon = 604 // polygon
	FieldTypeCircle  = 718 // circle

	// 网络地址类型
	FieldTypeCIDR    = 650 // cidr
	FieldTypeINET    = 869 // inet
	FieldTypeMacAddr = 829 // macaddr

	// JSON 类型
	FieldTypeJSON  = 114  // json
	FieldTypeJSONB = 3802 // jsonb

	// UUID
	FieldTypeUUID = 2950 // uuid

	// 数组类型示例
	FieldTypeInt4Array = 1007 // int4[]
)

type CustomStmt struct {
	id    int
	query string
}

func (s *CustomStmt) Close() error {
	log.Println("CustomStmt Close", s.id)
	return nil
}

func (s *CustomStmt) Exec(args []driver.Value) (driver.Result, error) {
	var res baseExecResult
	if err := query(context.Background(), s.id, s.query, args, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *CustomStmt) Query(args []driver.Value) (driver.Rows, error) {
	var res rawRows
	if err := query(context.Background(), s.id, s.query, args, &res); err != nil {
		return nil, err
	}
	res.init()
	return &res, nil
}

func (s *CustomStmt) NumInput() int {
	return -1
}
