package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"mo-gateway/src/config"
	"vilgo/vlog"
	"vilgo/vsql"
)

var (
	gblSql *MySqlServe
)

//
type MySqlServe struct{
	gblSql *vsql.NormalSqlDrive
}

func (sel *MySqlServe) Init(args ...interface{}) error {
	sel.gblSql = vsql.NewNormalSqlDrive(config.GetMySqlCnf())

	dbName := config.GetDbName()
	if err := sel.gblSql.Open(); err != nil {
		err = sel.createDatabase(sel.gblSql.GetDefDb(), dbName)
		if err != nil {
			vlog.LogE("init create database fail:%v", err)
			return err
		}
		err := sel.gblSql.Add(dbName)
		if err != nil {
			return err
		}
	}

	if dbName != "" {
		if err := sel.createAllTable(sel.gblSql.GetDb(dbName)); err != nil {
			vlog.LogE("init create table fail:%v", err)
			return err
		}
	}

	return nil
}

func (sel *MySqlServe) UnInit() {
	sel.gblSql.BreakOff()
}

// 创建数据库
func (*MySqlServe) createDatabase(db *sql.DB, dbName string) error {
	if db == nil {
		return errors.New("db is nil")
	}
	dbList, err := getDatabases(db)
	if err != nil {
		return err
	}
	if _, ok := dbList[dbName]; !ok {
		if _, err := db.Exec(fmt.Sprintf("create database %s;", dbName)); err != nil {
			return err
		}
	}
	return nil
}

// 创建数据库表
func (*MySqlServe) createAllTable(db *sql.DB) error {
	if db == nil {
		return nil
	}
	// 获取需要创建的表以及 创建表的 sql 语句
	tables := tableCreateSql()
	for k, tbn := range tables {
		// 获取现有的数据库表
		tbList, err := getTables(db)
		if err != nil {
			return err
		}

		// 判断是否存在数据库表，有的话就不用重新创建
		if _, ok := tbList[k]; ok {
			continue
		}

		// 数据库表不存在就创建
		if _, err := db.Exec(tbn); err != nil {
			return err
		}
	}
	return nil
}

func (sel *MySqlServe)GetDb()*sql.DB{
	return sel.gblSql.GetDb(config.GetDbName())
}

func Db() *sql.DB {
	return gblSql.GetDb()
}
