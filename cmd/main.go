package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sudoku/internal/domain"
	"sudoku/internal/user/delivery/httpdeliver"
	"sudoku/internal/user/repository/mysqlHandler"
	"sudoku/internal/user/usecase"
)

func main() {
	dbUser := "root"
	dbPassword := "97216017"
	dbName := "sudoku"
	dsn := dbUser + ":" + dbPassword + "@tcp(127.0.0.1:3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	Db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Connecting to database failed")
	}
	_ = Db.AutoMigrate(&domain.User{})
	_ = Db.AutoMigrate(&domain.Board{})
	r := echo.New()
	ur := mysqlHandler.NewMysqlUserRepository(Db)
	uu := usecase.NewUserUsecase(ur)
	httpdeliver.NewUserHandler(r, uu)
}
