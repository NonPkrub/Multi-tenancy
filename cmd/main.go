package main

import (
	"fmt"
	"go-multi-tenancy/internals/core/services"
	"go-multi-tenancy/internals/handlers"
	"go-multi-tenancy/internals/repositories"
	"go-multi-tenancy/internals/server"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	initConfig()
	db := initDatabase()

	companyRepository := repositories.NewCompanyRepository(db)
	companyService := services.NewCompanyService(companyRepository)
	companyHandler := handlers.NewCompanyHandler(companyService)

	manageRepository := repositories.NewManageRepository(db)
	manageService := services.NewManageService(manageRepository)
	manageHandler := handlers.NewManageHandler(manageService)

	httpServer := server.NewServer(companyHandler, manageHandler)

	httpServer.Initialize()
}

func initDatabase() *sqlx.DB {
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable&timezone=Asia/Bangkok",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetInt("db.port"),
		viper.GetString("db.database"),
	)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("..")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
