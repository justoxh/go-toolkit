package mysql

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type MysqlConfig struct {
	Hostname    string
	Port        string
	User        string
	Password    string
	DbName      string
	TablePrefix string
}

// MysqlService implements mysql service
type MysqlService struct {
	config MysqlConfig
	DB     *sql.DB
}

// NewMysqlService create new mysql service
func NewMysqlService(conf MysqlConfig) *MysqlService {
	service := &MysqlService{}
	url := conf.User + ":" + conf.Password + "@tcp(" + conf.Hostname + ":" + conf.Port + ")/" + conf.DbName + "?charset=utf8&parseTime=True"
	db, err := sql.Open("mysql", url)
	if err != nil {
		// todo : add log
		return nil
	}
	service.DB = db
	service.config = conf
	return service
}
