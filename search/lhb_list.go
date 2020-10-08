package search

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const formatLayout = "2006-01-02"

type (
	lhbList struct {
		startDate string
		endDate   string
		data      list.List
		pageSize  int
		wg        *sync.WaitGroup
		db        *gorm.DB
	}

	Result struct {
		data         list.List
		dateInterval string
		currPage     int
		pageSize     int
	}

	Data struct {
		TableName      string
		ConsumeMSecond int
		Data           []string
		FieldName      string
		SplitSymbol    string
		TotalPage      int
	}

	LHBGGDRTJ struct {
		Data    []Data
		Message string
		Status  int
	}
)

var u = &Utils{}

func (ll *lhbList) fetch(page int, wg *sync.WaitGroup) {
	params := u.GetParams(ll.startDate, ll.endDate, page)
	body, err := u.ListRequest(params)
	if err != nil {
		log.Printf("producer:ERROR:ListRequest: url HTTP Response error:%s\n", err)
		wg.Done()
		return
	}
	data, err := ll.parser(body)
	if err != nil {
		log.Printf("producer:ERROR:parser: %s\n", err)
		wg.Done()
		return
	}
	ll.data.PushBackList(&data)
	wg.Done()
}

func (ll *lhbList) save() {
	filepath := fmt.Sprintf("out/%s_%s.txt", ll.startDate, ll.endDate)
	f, _ := os.Create(filepath)
	defer f.Close()
	for d := ll.data.Front(); d != nil; d = d.Next() {
		_, _ = f.WriteString(fmt.Sprintf("%s\n", d.Value))
		llist := &LHBList{}
		_ = mapstructure.Decode(d.Value, llist)
		ll.db.Create(&EastMoney{llist, &LHBInfo{}})
	}
}

func (ll *lhbList) parser(body string) (list.List, error) {
	var result list.List
	lhb := &LHBGGDRTJ{}
	if err := json.Unmarshal([]byte(body), lhb); err != nil {
		return result, err
	}
	data := lhb.Data[0]
	fieldNames := strings.Split(data.FieldName, ",")
	for _, item := range data.Data {
		r := make(map[string]interface{})
		for i, field := range strings.Split(item, data.SplitSymbol) {
			r[fieldNames[i]] = field
		}
		result.PushBack(r)
	}
	ll.pageSize = data.TotalPage
	return result, nil
}

func (ll *lhbList) getPages() int {
	params := u.GetParams(ll.startDate, ll.endDate, 1)
	body, _ := u.ListRequest(params)
	_, err := ll.parser(body)
	if err == nil {
		return ll.pageSize
	}
	log.Printf("search:getPages:ERROR: %s\n", err)
	return -1
}

func (ll *lhbList) producer() {
	pages := ll.getPages()
	wg := &sync.WaitGroup{}
	wg.Add(pages)
	for page := 1; page <= pages; page++ {
		go ll.fetch(page, wg)
		time.Sleep(time.Millisecond)
	}
	wg.Wait()
	ll.save()
	ll.wg.Done()
}

func LHBListProducer(db *gorm.DB) {
	startDate := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Now()
	wg := &sync.WaitGroup{}
	for {
		// 以60天步进
		ed := startDate.Add(time.Hour * 24 * 60)
		if !endDate.After(ed) {
			break
		}
		sdStr := startDate.Format(formatLayout)
		edStr := ed.Format(formatLayout)

		wg.Add(1)
		ll := &lhbList{
			startDate: sdStr,
			endDate:   edStr,
			wg:        wg,
			db:        db,
		}

		go ll.producer()
		time.Sleep(time.Millisecond * 200)
		startDate = ed
	}
	wg.Wait()
	fmt.Println("Fetch Done. ")
}
