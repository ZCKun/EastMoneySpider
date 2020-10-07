package search

import (
	"fmt"
	"gorm.io/gorm"
)

type (
		info struct {

		}
)

func (i *info) Do(db *gorm.DB) {
	var result []map[string]interface{}
	db.Model(&EasyMoney{}).Select("s_code", "s_name").Find(&result)
	fmt.Println(result)
}
