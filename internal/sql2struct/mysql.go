package sql2struct

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DBModel struct {
	DBEngine *sql.DB
	DBInfo   *DBInfo
}

type DBInfo struct {
	DBType   string
	Host     string
	UserName string
	Password string
	Charset  string
}

type Column struct {
	ColumnName    string
	DataType      string
	IsNullable    string
	ColumnKey     string
	ColumnType    string
	ColumnComment string
	TableName     string
}

type TableColumnAndTableName struct {
	TableName    string
	TableColumns []*Column
}

func NewDBModel(info *DBInfo) *DBModel {
	return &DBModel{DBInfo: info}
}

func (m *DBModel) Connect() error {
	var err error
	s := "%s:%s@tcp(%s)/information_schema?" +
		"charset=%s&parseTime=True&loc=Local"

	dsn := fmt.Sprintf(
		s,
		m.DBInfo.UserName,
		m.DBInfo.Password,
		m.DBInfo.Host,
		m.DBInfo.Charset,
	)

	m.DBEngine, err = sql.Open(m.DBInfo.DBType, dsn)
	if err != nil {
		return err

	}
	return nil
}

func (m *DBModel) GetColumns(dbName, tableName string) ([]*Column, error) {
	query := "SELECT COLUMN_NAME,DATA_TYPE,COLUMN_KEY,IS_NULLABLE,COLUMN_TYPE,COLUMN_COMMENT " +
		" FROM COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? "

	rows, err := m.DBEngine.Query(query, dbName, tableName)
	if err != nil {
		return nil, err

	}
	if rows == nil {
		return nil, errors.New("没有数据")

	}

	defer rows.Close()
	var columns []*Column
	for rows.Next() {
		var column Column
		err := rows.Scan(&column.ColumnName, &column.DataType, &column.ColumnKey, &column.IsNullable, &column.ColumnType, &column.ColumnComment)
		if err != nil {
			return nil, err
		}
		column.TableName = tableName

		columns = append(columns, &column)
	}

	return columns, nil

}

var DBTypeToStructType = map[string]string{
	"int":       "int",
	"tinyint":   "int",
	"bigint":    "uint64",
	"bit":       "int",
	"decimal":   "float64",
	"bool":      "bool",
	"enum":      "string",
	"timestamp": "*time.Time",
	"set":       "string",
	"varchar":   "string",
	"text":      "string",
}
