package db

import (
	"database/sql"

	"time"

	"sync"
	"../conf"
)

var (
	conn         *sql.DB
	newMysqlLock = new(sync.Mutex)
)

// 从数据库连接池获得连接
func Conn() *sql.DB {
	if conn == nil {
		newMysqlLock.Lock()
		defer newMysqlLock.Unlock()
		if conn == nil {
			edmCoreConf := conf.NewConf()
			mysql := edmCoreConf.MySql
			var err error
			conn, err = sql.Open("mysql", mysql.Url)
			if err != nil {
				panic(err.Error())
			}
			conn.SetMaxIdleConns(mysql.MaxIdleConns)
			conn.SetMaxOpenConns(mysql.MaxOpenConns)
			conn.SetConnMaxLifetime(time.Duration(mysql.ConnMaxLifetime) * time.Second)
		}
	}
	return conn
}

// 程序停止时关闭连接池
func Close() {
	if conn != nil {
		conn.Close()
	}
}
