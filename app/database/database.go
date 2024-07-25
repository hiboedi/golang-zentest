package database

import (
	"fmt"
	"log"

	"zen-test/app/helpers"
	"zen-test/app/web/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBDriver   string
}

func InitializeDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error on loading .env file")
	}
	var dbConfig DBConfig

	dbConfig.DBName = helpers.GetEnv("DATABASE_NAME", "exam")
	dbConfig.DBPort = helpers.GetEnv("DATABASE_PORT", "5432")
	dbConfig.DBUser = helpers.GetEnv("DATABASE_USER", "boedi")
	dbConfig.DBPassword = helpers.GetEnv("DATABASE_PASSWORD", "")
	dbConfig.DBHost = helpers.GetEnv("DATABASE_HOST", "localhost")
	dbConfig.DBDriver = helpers.GetEnv("DATABASE_DRIVER", "postgres")

	if dbConfig.DBDriver == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
		helpers.PanicIfError(err)
		return db
	} else {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBPort, dbConfig.DBName)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

		if err != nil {
			panic("Failed on connecting to the database server")
		}
		return db
	}
}

func DBMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		// &models.Category{},
		&models.Order{},
		&models.Image{},
		&models.Product{},
		&models.OrderItem{},
	)
	helpers.PanicIfError(err)

	fmt.Println("Db migration success")
}
