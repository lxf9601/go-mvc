// 系统配置
// 对应sys_conf.json的配置结构，根据需要增改，字段保持与文件配置项一一对应即可
package conf

import (
	"sync"

	"path/filepath"

	"github.com/jinzhu/configor"
	"cparrow.com/go-mvc/util"
)

// 系统配置节点
type SysConf struct {
	MySql              *MySqlConf // mysql配置
	Redis              *RedisConf // redis配置
	AppPath            string     // 应用当前物理路径
	LogLevel           string     // 系统日志级别
	ServerURL          string     // WEB服务器地址
}

// Mysql配置项
type MySqlConf struct {
	Url             string // 连接字符串
	MaxIdleConns    int    // 最大空闲连接数
	MaxOpenConns    int    // 最大打开连接数
	ConnMaxLifetime int    // 连接生命周期
}

// Redis配置项
type RedisConf struct {
	Addr      string // redis服务地址
	Password  string // 密码
	Db        int    // 数据库0
	KeyPrefix string // 缓存key值前缀
}

// 配置实例
var sysConf *SysConf = nil

// 单例构建锁
var sysConfLock = new(sync.Mutex)

// 配置文件路径
var confFilePath = "conf/sys_conf.json"

// 创建配置（单例）
func NewConf() *SysConf {
	if sysConf == nil {
		sysConfLock.Lock()
		defer sysConfLock.Unlock()
		if sysConf == nil {
			sysConf = new(SysConf)
			mysql := new(MySqlConf)
			sysConf.MySql = mysql
			sysConf.AppPath, _ = util.GetAppPath()
			confPath := sysConf.AppPath + confFilePath
			exists, _ := util.PathExists(confPath)
			if !exists {
				path, _ := filepath.Abs("./")
				sysConf.AppPath = path + "/"
				confPath = sysConf.AppPath + confFilePath
				exists, _ := util.PathExists(confPath)
				if !exists {
					path, _ := filepath.Abs("../")
					sysConf.AppPath = path + "/"
					confPath = sysConf.AppPath + confFilePath

				}
			}
			configor.Load(&sysConf, confPath)
		}
	}
	return sysConf
}
