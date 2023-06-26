package sqlDiff_mysql8_0_32

import (
	"database/sql"
	"github.com/sirupsen/logrus"
)

func getColumnsComments(db *sql.DB, tableName string) (ans map[string]string, err error) {
	rows, err := db.Query("SELECT COLUMN_NAME, COLUMN_COMMENT FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = ?", tableName)
	if err != nil {
		return nil, err
	}
	ans = make(map[string]string)
	defer rows.Close()
	for rows.Next() {
		var field string
		var comment string
		err = rows.Scan(&field, &comment)
		if err != nil {
			logrus.WithError(err).Errorln("fail in scan result from mysql")
			continue
		}
		if len(comment) != 0 {
			ans[field] = comment
		}
	}
	return
}

// 获取表的列信息
func getColumns(db *sql.DB, tableName string) (columns []*Column, err error) {
	rows, err := db.Query("SHOW COLUMNS FROM " + tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		column := new(Column)
		err := rows.Scan(
			&column.Field,
			&column.Type,
			&column.Null,
			&column.Key,
			&column.Default,
			&column.Extra,
		)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	return columns, nil
}

func getIndexes(db *sql.DB, tbName string) (ans []Index, err error) {
	rows, err := db.Query("SHOW INDEX FROM " + tbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var index Index
		err := rows.Scan(
			&index.Table,
			&index.NonUnique,
			&index.KeyName,
			&index.SeqInIndex,
			&index.ColumnName,
			&index.Collation,
			&index.Cardinality,
			&index.SubPart,
			&index.Packed,
			&index.Null,
			&index.IndexType,
			&index.Comment,
			&index.IndexComment,
			&index.Visible,
			&index.Expression,
		)
		if err != nil {
			logrus.WithError(err).Errorln("fail to get index from db")
			continue
		}
		flag := false
		for i, idx := range ans {
			if index.KeyName == idx.KeyName {
				idx.ColumnName += "," + index.ColumnName
				ans[i] = idx
				flag = true
				break
			}
		}
		if flag {
			continue
		}

		// 将当前行的数据添加到切片中
		ans = append(ans, index)
	}

	return ans, nil
}

func getTables(db *sql.DB) (tableNames []string, err error) {
	rows, err := db.Query("SHOW TABLES")
	for rows.Next() {
		var tbName string
		err = rows.Scan(&tbName)
		if err != nil {
			logrus.WithError(err).Errorln("fail to get tables from db")
			continue
		}
		tableNames = append(tableNames, tbName)
	}
	return
}

func getCreateTables() (ans []string) {
	for _, table := range sourceTables {
		if !tableExists(targetTables, table) {
			rows, err := sourceDB.Query("SHOW CREATE TABLE " + table)
			if err != nil {
				logrus.WithError(err).Errorln("fail to show create table " + table)
			}
			for rows.Next() {
				var tableName string
				var createTableSql string
				err := rows.Scan(&tableName, &createTableSql)
				if err != nil {
					logrus.WithError(err).Errorln("fail to get create table sql")
					continue
				}

				ans = append(ans, createTableSql)
			}
		}
	}
	return
}
