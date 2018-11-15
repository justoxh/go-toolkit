package mysqldao

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // init mysql
)

// MysqlConifg implements config
type MysqlConifg struct {
	Hostname    string
	Port        string
	User        string
	Password    string
	DbName      string
	TablePrefix string

	MaxOpenConnections int
	MaxIdleConnections int
	ConnMaxLifetime    int // unit second
	Debug              bool
	Local       string
}

// RdsService
type RdsService struct {
	config MysqlConifg
	DB     *gorm.DB
}

// NewRdsService create new mysql dao service
func NewRdsService(config MysqlConifg) (*RdsService, error) {
	impl := &RdsService{}
	impl.config = config

	// "root@tcp(127.0.0.1:3306)/s3d?charset=utf8"
	password := config.Password
	if password != "" {
		password = fmt.Sprintf(":%s", password)
	}
	if config.Local =="" {
		url := fmt.Sprintf("%s%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True&", config.User, password, "tcp", config.Hostname, config.Port,
		config.DbName)
	}else{
		url := fmt.Sprintf("%s%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s", config.User, password, "tcp", config.Hostname, config.Port,
		config.DbName,config.Local)
	}
	

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return config.TablePrefix + defaultTableName
	}
	db, err := gorm.Open("mysql", url)
	if err != nil {
		log.Fatalf("mysql connection error:%s", err.Error())
		return nil, err
	}

	maxOpenConns := config.MaxOpenConnections
	if maxOpenConns < 5 {
		maxOpenConns = 5
	}

	db.DB().SetMaxOpenConns(maxOpenConns)

	maxIdleConns := config.MaxIdleConnections
	if maxIdleConns < 1 {
		maxIdleConns = 1
	}
	db.DB().SetMaxIdleConns(maxIdleConns)

	connMaxLifeTime := config.ConnMaxLifetime
	if connMaxLifeTime < 30 {
		connMaxLifeTime = 30
	}
	db.DB().SetConnMaxLifetime(time.Duration(connMaxLifeTime) * time.Second)

	db.LogMode(config.Debug)
	err = db.DB().Ping()
	if err != nil {
		panic(fmt.Sprintf("init mysql db err: %v", err))
	}

	impl.DB = db

	return impl, nil
}

// RegistTable create table for given object
func (s *RdsService) RegistTable(t interface{}) {
	if ok := s.DB.HasTable(t); !ok {
		if err := s.DB.CreateTable(t).Error; err != nil {
			log.Fatalf("create mysql table error:%s", err.Error())
		}
	}
	var tab []interface{}
	s.DB.AutoMigrate(append(tab, t))
}

// RegistTables create tables for given object
func (s *RdsService) RegistTables(tables []interface{}) {
	for _, t := range tables {
		if ok := s.DB.HasTable(t); !ok {
			if err := s.DB.CreateTable(t).Error; err != nil {
				log.Fatalf("create mysql table error:%s", err.Error())
			}
		}
	}

	// auto migrate to keep schema update to date
	// AutoMigrate will ONLY create tables, missing columns and missing indexes,
	// and WON'T change existing column's type or delete unused columns to protect your data
	s.DB.AutoMigrate(tables...)
}

// Close db
func (s *RdsService) Close() error {
	return s.DB.Close()
}

// Add single item
func (s *RdsService) Add(item interface{}) error {
	return s.DB.Create(item).Error
}

// Del single item
func (s *RdsService) Del(item interface{}) error {
	return s.DB.Delete(item).Error
}

// Save single item
func (s *RdsService) Save(item interface{}) error {
	return s.DB.Save(item).Error
}
