package main

import (
	"EastMoneySpider/search"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func initDB() (*gorm.DB, error) {
	dsn := "root:aaaa@tcp(127.0.0.1:3307)/eastmoney?charset=utf8mb4"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		_ = fmt.Errorf("Connect mysql failed: %s\n", err)
		return nil, err
	} else {
		em := search.EastMoney{}
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
    // search.LHBListProducer(db)
    search.Search(db)
	//
}

//func main() {
//	u := utils.ReasonContains("日换手率达到24.07%", "2020-03-27 星期五  类型：换手率达20%的证券")
//	if u {
//		println("y")
//	} else {
//		println("x")
//	}
//}