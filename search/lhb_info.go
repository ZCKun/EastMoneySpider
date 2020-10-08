package search

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
)

type (
	info struct {
		db *gorm.DB
	}
)

func (i *info) fetch(m map[string]interface{}) {
	date := m["tdate"].(string)
	code := m["s_code"].(string)
	body, err := u.InfoRequest(date, code)
	if err != nil {
		log.Printf("lhb_info:fetch:Request Info Error: %s\n", err)
		return
	}
	i.parser(body, m)
}

func (i *info) parser(body string, m map[string]interface{}) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.Printf("lhb_info:parser:parse body to html error: %s\n", err)
	}
	lInfo := &LHBInfo{}
	doc.Find(".content-sepe table").Each(func(i int, s *goquery.Selection) {
		tableClass, _ := s.Attr("class")
		fmt.Println(tableClass)
		if strings.Contains(tableClass, "stock-detail-tab") {
			topBuyParser(s, lInfo)
		} else {
			topSellParser(s, lInfo)
		}
	})
	fmt.Println(lInfo)
	code := m["s_code"].(string)
	date := m["tdate"].(string)
	name := m["s_name"].(string)
	ctypeDes := m["ctypedes"].(string)
	dp := m["dp"].(string)
	i.db.Model(&EasyMoney{}).Where("s_code = ? AND s_name = ? AND tdate = ? AND ctypedes = ? AND dp = ?",
		code, name, date, ctypeDes, dp).Updates(EasyMoney{nil, lInfo})
}

func (i *info) Do() {
	var result []map[string]interface{}
	i.db.Model(&EasyMoney{}).Select("s_code", "s_name", "tdate", "ctypedes", "dp").Find(&result)
	for _, m := range result {
		i.fetch(m)
	}
}

func topSellParser(s *goquery.Selection, info *LHBInfo) {
	for _, tr := range s.Find("tbody tr").Nodes {
		topSell := &DealInfo{}
		// 过滤掉最后的总计和
		if len(tr.Attr) > 0 {
			continue
		}
		n := goquery.NewDocumentFromNode(tr)
		for i, td := range n.Find("td").Nodes {
			tdDoc := goquery.NewDocumentFromNode(td)
			d := tdDoc.Has("div")
			if len(d.Nodes) > 0 {
				findSCName(td, topSell)
			} else {
				findDealInfo(i, tdDoc, topSell)
			}
		}
		ts, _ := json.Marshal(topSell)
		info.TopSell = append(info.TopSell, string(ts))
	}
}

func topBuyParser(s *goquery.Selection, info *LHBInfo) {
	for _, tr := range s.Find("tbody tr").Nodes {
		topBuy := &DealInfo{}
		// 过滤掉最后的总计和
		if len(tr.Attr) > 0 {
			continue
		}
		n := goquery.NewDocumentFromNode(tr)
		for i, td := range n.Find("td").Nodes {
			tdDoc := goquery.NewDocumentFromNode(td)
			d := tdDoc.Has("div")
			if len(d.Nodes) > 0 {
				findSCName(td, topBuy)
			} else {
				findDealInfo(i, tdDoc, topBuy)
			}
		}
		tb, _ := json.Marshal(topBuy)
		info.TopBuy = append(info.TopBuy, string(tb))
	}
}

func findDealInfo(i int, tdDoc *goquery.Document, topAny *DealInfo) {
	switch i {
	case 2:
		//  买入金额(万)
		f, _ := strconv.ParseFloat(tdDoc.Text(), 64)
		topAny.BuyAmount = f * 10000
		break
	case 3:
		// 占总成交比例
		topAny.BuyAmountPropTotalTran = tdDoc.Text()
		break
	case 4:
		// 卖出金额(万)
		f, _ := strconv.ParseFloat(tdDoc.Text(), 64)
		topAny.SellAmount = f * 10000
		break
	case 5:
		// 占总成交比例
		topAny.SellAmountPropTotalTran = tdDoc.Text()
		break
	case 6:
		// 净额(万)
		f, _ := strconv.ParseFloat(tdDoc.Text(), 64)
		topAny.NetAmount = f * 10000
		break
	}
}

func findSCName(td *html.Node, topAny *DealInfo) {
	for _, a := range goquery.NewDocumentFromNode(td).Find(".sc-name a").Nodes {
		if scName := goquery.NewDocumentFromNode(a).Text(); scName != "" {
			topAny.SCName = scName
		}
		for _, _a := range a.Attr {
			if _a.Key == "href" && _a.Val != "" && strings.ContainsAny(_a.Val, "lhb/yyb") {
				topAny.Href = _a.Val
			}
		}
	}
}
