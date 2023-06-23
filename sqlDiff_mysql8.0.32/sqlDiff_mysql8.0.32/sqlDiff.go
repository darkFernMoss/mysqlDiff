package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"
)

var targetDBUrl = flag.String("targetDB", "", "url of target database, example: 'username:password@tcp(ip:port)/myDB'")
var sourceDBUrl = flag.String("sourceDB", "", "url of source database, example: 'username:password@tcp(ip:port)/myDB'")

var sourceTables []string
var targetTables []string
var sourceDB *sql.DB
var targetDB *sql.DB

func main() {
	initDB()
	printUpdateSql()
}

func initDB() {
	flag.Parse()
	args := os.Args
	if len(args) != 5 {
		fmt.Println("please use ./sqlDiff -h for help")
		os.Exit(1)
	}
	if targetDBUrl == nil || len(*targetDBUrl) == 0 {
		logrus.Fatal("target database path must be specified correctly")
	}
	if targetDBUrl == nil || len(*sourceDBUrl) == 0 {
		logrus.Fatal("source database path must be specified correctly")
	}

	// 连接数据库
	var err error
	targetDB, err = sql.Open("mysql", *targetDBUrl)
	if err != nil {
		logrus.WithError(err).Fatal("the url of targetDB is wrong")
	}
	err = targetDB.Ping()
	if err != nil {
		logrus.WithError(err).Fatal("cannot connect to targetDB")
	}
	sourceDB, err = sql.Open("mysql", *sourceDBUrl)
	if err != nil {
		logrus.WithError(err).Fatal("the url of sourceDB is wrong")
	}
	err = sourceDB.Ping()
	if err != nil {
		logrus.WithError(err).Fatal("cannot connect to sourceDB")
	}

	targetTables, err = getTables(targetDB)
	if err != nil {
		logrus.WithError(err).Fatalf("fail to get all tables from %s", *targetDBUrl)
	}
	sourceTables, err = getTables(sourceDB)
	if err != nil {
		logrus.WithError(err).Fatalf("fail to get all tables from %s", *sourceDBUrl)
	}
}

func printUpdateSql() {
	printDbDiff()
	for _, ttb := range targetTables {
		for _, stb := range sourceTables {
			if ttb == stb {
				printSqlDiff(ttb)
			}
		}
	}
}

func printDbDiff() {
	tBaseName := strings.Split(*targetDBUrl, "/")[1]
	sBaseName := strings.Split(*sourceDBUrl, "/")[1]

	var dbDiff []string
	for _, table := range targetTables {
		if !tableExists(sourceTables, table) {
			dbDiff = append(dbDiff, fmt.Sprintf("DROP TABLE %s;", table))
		}
	}
	createTableSql := getCreateTables()
	dbDiff = append(dbDiff, createTableSql...)
	if len(dbDiff) != 0 {
		fmt.Println()
		fmt.Println()
		fmt.Printf("#Table changes from database %s to database %s:\n", tBaseName, sBaseName)
		for _, sqlStr := range dbDiff {
			fmt.Println(sqlStr + ";")
		}
	}
	fmt.Println()
	fmt.Println()
}

func printSqlDiff(tableName string) {

	// 获取源表的列信息
	sourceColumns, err := getColumns(sourceDB, tableName)
	if err != nil {
		log.Fatal(err)
	}
	// 获取目标表的列信息
	targetColumns, err := getColumns(targetDB, tableName)
	if err != nil {
		log.Fatal(err)
	}

	sComments, err := getColumnsComments(sourceDB, tableName)
	if err != nil {
		logrus.WithError(err).Fatal("fail to get sComments")
	}

	tComments, err := getColumnsComments(targetDB, tableName)
	if err != nil {
		logrus.WithError(err).Fatal("fail to get sComments")
	}

	tIndexes, err := getIndexes(targetDB, tableName)
	if err != nil {
		logrus.WithError(err).Fatal("fail to get targetIndexes")
	}
	sIndexes, err := getIndexes(sourceDB, tableName)
	if err != nil {
		logrus.WithError(err).Fatal("fail to get sourceIndexes")
	}

	// 比较源表和目标表的列差异
	columnDiff := make([]string, 0, 100)
	indexDiff := make([]string, 0, 10)
	addColumn(sourceColumns, targetColumns, sComments, &columnDiff)
	dropColumn(targetColumns, sourceColumns, &columnDiff)
	modifyColumn(targetColumns, sourceColumns, sComments, tComments, &columnDiff)

	//比较索引差异
	addIndex(sIndexes, tIndexes, &indexDiff)
	dropIndex(sIndexes, tIndexes, &indexDiff)

	// 生成表变更的 DDL 语句
	printResult(columnDiff, indexDiff, tableName)
}

func printResult(columnDiff []string, indexDiff []string, tableName string) {
	if len(columnDiff) > 0 {
		ddlStatement := fmt.Sprintf("ALTER TABLE %s\n%s", tableName, joinStrings(columnDiff, ",\n"))
		fmt.Printf("#The sql statement corresponding to the change of table %s ：\n%s;\n", tableName, ddlStatement)
		fmt.Println()
		fmt.Println()
		fmt.Printf("#Index update statement for table %s：\n", tableName)
		for _, idxStatement := range indexDiff {
			fmt.Println(idxStatement)
		}
		fmt.Println()
		fmt.Println()
	}
}

func addIndex(sIndexes []Index, tIndexes []Index, indexDiff *[]string) {
	for _, sIdx := range sIndexes {
		if !indexExists(&tIndexes, sIdx.KeyName) {
			uniqueStr := " "
			if sIdx.NonUnique == 0 {
				uniqueStr = " UNIQUE "
			}
			sqlStr := fmt.Sprintf("CREATE%sINDEX %s ON %s(%s);", uniqueStr, sIdx.KeyName, sIdx.Table, sIdx.ColumnName)
			*indexDiff = append(*indexDiff, sqlStr)
		}
	}
}

func dropIndex(sIndexes []Index, tIndexes []Index, indexDiff *[]string) {
	for _, tIdx := range tIndexes {
		if !indexExists(&sIndexes, tIdx.KeyName) {
			sqlStr := fmt.Sprintf("DROP INDEX %s;", tIdx.KeyName)
			*indexDiff = append(*indexDiff, sqlStr)
		}
	}
}

func modifyColumn(targetColumns []*Column, sourceColumns []*Column, sComments map[string]string, tComments map[string]string, columnDiff *[]string) {
	for _, tclm := range targetColumns {
		for _, sclm := range sourceColumns {
			if tclm.Field == sclm.Field {
				if tclm.Type != sclm.Type || tclm.Null != sclm.Null || tclm.Key != sclm.Key ||
					!defaultEqual(tclm.Default, sclm.Default) || sclm.Extra != tclm.Extra || !commentEqual(sComments, tComments, sclm.Field, tclm.Field) {
					sqlStr := fmt.Sprintf("MODIFY COLUMN %s %s", sclm.Field, sclm.Type)
					if sclm.Null == "NO" {
						sqlStr += " NOT NULL"
					}
					if len(sclm.Key) != 0 {
						switch sclm.Key {
						case "PRI":
							sqlStr += " PRIMARY KEY"
						}
					}
					if sclm.Default != nil {
						bts := sclm.Default.([]byte)
						if len(bts) > 0 {
							sqlStr += " DEFAULT " + string(bts)
						}
					}
					if len(sclm.Extra) != 0 && sclm.Extra == "AUTO_INCREMENT" {
						sqlStr += " " + sclm.Extra
					}
					if cmt, ok := sComments[sclm.Field]; ok {
						sqlStr += " COMMENT " + "'" + cmt + "'"
					}
					*columnDiff = append(*columnDiff, sqlStr)
				}

			}
		}
	}

}

func dropColumn(targetColumns []*Column, sourceColumns []*Column, columnDiff *[]string) {
	for _, targetColumn := range targetColumns {
		columnName := targetColumn.Field

		if !columnExists(sourceColumns, columnName) {
			*columnDiff = append(*columnDiff, fmt.Sprintf("DROP COLUMN %s", columnName))
		}
	}
}

func addColumn(sourceColumns []*Column, targetColumns []*Column, sComments map[string]string, columnDiff *[]string) {
	for _, clm := range sourceColumns {
		columnName := clm.Field
		columnType := clm.Type

		if !columnExists(targetColumns, columnName) {
			addColumn := fmt.Sprintf("ADD COLUMN %s %s", columnName, columnType)
			if clm.Null == "NO" {
				addColumn += " NOT NULL"
			}
			if len(clm.Key) != 0 {
				switch clm.Key {
				case "PRI":
					addColumn += " PRIMARY KEY"
				}
			}
			if clm.Default != nil {
				bts := clm.Default.([]byte)
				addColumn += " DEFAULT" + string(bts)
			}
			if len(clm.Extra) != 0 && clm.Extra == "AUTO_INCREMENT" {
				addColumn += " " + clm.Extra
			}
			if cmt, ok := sComments[columnName]; ok {
				addColumn += " COMMENT " + "'" + cmt + "'"
			}
			*columnDiff = append(*columnDiff, addColumn)
		}
	}
}
