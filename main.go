package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func initDB() *gorm.DB{
	db, err := gorm.Open("postgres", os.Getenv("DB_URI"))
	if err != nil {
		panic(err)
	}

	return db
}


func main() {
	db := initDB()
	defer db.Close()

	fmt.Printf("Connected to database: \n%s\n", os.Getenv("DB_URI"))
}
