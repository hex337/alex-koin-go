package Config

import (
  // "log"
	"fmt"

	// "gorm.io/gorm/logger"
	"gorm.io/gorm"
)

var DB *gorm.DB

type DBConfig struct {
	Host	   string
	Port	   int
	User	   string
	Password string
	DBName   string
}

func BuildDBConfig() *DBConfig {
	dbConfig := DBConfig{
  Host:     "localhost",
  Port:     5432,
  User:     "pbaker",
  Password: "foo",
  DBName:   "akc",
  }
  return &dbConfig
}

func DBURL(dbConfig *DBConfig) string {
  return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.Port,
	)
}
