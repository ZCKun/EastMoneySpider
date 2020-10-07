package main

import (
	"EasyMoneySpider/search"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func initDB() (*gorm.DB, error) {
	dsn := "root:aaa111@tcp(0.0.0.0:3306)/eastmoney?charset=utf8mb4"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		_ = fmt.Errorf("Connect mysql failed: %s\n", err)
		return nil, err
	} else {
		em := search.EasyMoney{}
		if !db.Migrator().HasTable(em) {
			_ = db.Migrator().CreateTable(em)
		}
		return db, nil
	}
}

func init() {
	dt := time.Now().Format("20060102150405")
	f, _ := os.Create(fmt.Sprintf("main_%s.log", dt))
	log.SetOutput(f)
}

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal("init db has found error.")
	}
	// search.Do(db)
	search.Search(db)
}