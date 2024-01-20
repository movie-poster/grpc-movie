package config

import (
	"fmt"
	"template-grpc/internal/domain/entity"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	err error
)

type Database struct {
	*gorm.DB
}

// SetupDB opens a database and saves the reference to `Database` struct.
func setupDB(configuration *Configuration) {
	var db = DB

	driver := configuration.Database.Driver
	database := configuration.Database.Dbname
	username := configuration.Database.Username
	password := configuration.Database.Password
	host := configuration.Database.Host
	port := configuration.Database.Port

	if driver == "sqlite" { // SQLITE
		db, err = gorm.Open(sqlite.Open(database+".db"), &gorm.Config{})
		if err != nil {
			fmt.Println("db err: ", err)
		}
	} else if driver == "postgres" { // POSTGRES
		db, err = gorm.Open(postgres.Open("host="+host+" port="+port+" user="+username+" dbname="+database+"  sslmode=disable password="+password), &gorm.Config{})
		if err != nil {
			fmt.Println("db err: ", err)
		}
	} else if driver == "mysql" { // MYSQL
		db, err = gorm.Open(mysql.Open(username+":"+password+"@tcp("+host+":"+port+")/"+database+"?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		})
		if err != nil {
			fmt.Println("db err: ", err)
		}
	}

	sqlDB, _ := db.DB()
	// Change this to true if you want to see SQL queries

	db.Logger.LogMode(logger.Info)
	sqlDB.SetMaxIdleConns(configuration.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(configuration.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(configuration.Database.MaxLifetime) * time.Second)
	DB = db
	migration()
}

// Auto migrate project models
func migration() {
	DB.AutoMigrate(&entity.User{})

}

func GetDB() *gorm.DB {
	return DB
}
