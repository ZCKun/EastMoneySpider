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
	"sync"
	"time"
)

type (
	info struct {
		db *gorm.DB
	}
)

func (i *info) fetch(em EastMoney, wg *sync.WaitGroup) {
	date := em.Tdate
	code := em.SCode
	body, err := u.InfoRequest(date, code)
	if err != nil {
		log.Printf("lhb_info:fetch:Request Info url %s Error: %s\n", body, err)
		return
	}
	i.parser(body, em)
	wg.Done()
}

func (i *info) parser(body string, em EastMoney) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.Printf("lhb_info:parser:parse body to html error: %s\n", err)
	}

	contentBox := doc.Find(".content-box")
	// 上榜原因明细
	reasons := i.reasonHandle(contentBox)
	lInfos := i.infoHandle(contentBox)

	// update to sql
	code := em.SCode
	date := em.Tdate
	name := em.SName
	ctypeDes := em.Ctypedes
	dp := em.DP

	var lInfo *LHBInfo
	for i, _ := range reasons {
		lInfo = &lInfos[i]
	}
	//for i, reason := range reasons {
	//	if utils.ReasonContains(ctypeDes, reason) {
	//		lInfo = &lInfos[i]
	//	}
	//}

	i.db.Model(&EastMoney{}).Where("s_code = ? AND s_name = ? AND tdate = ? AND ctypedes = ? AND dp = ?",
		code, name, date, ctypeDes, dp).Updates(EastMoney{nil, lInfo})
}

func (i *info) infoHandle(contentBox *goquery.Selection) []LHBInfo {
	var lInfos []LHBInfo
	contentBox.Find(".content-sepe").Each(func(i int, s *goquery.Selection) {
		lInfo := LHBInfo{}
		s.Find("table").Each(func(i int, s *goquery.Selection) {
			tableClass, _ := s.Attr("class")
			if strings.Contains(tableClass, "stock-detail-tab") {
				tb, totalBuyAmount := topBuyParser(s)
				lInfo.TopBuy = fmt.Sprintf("[%s]", strings.Join(tb, ","))
				lInfo.TotalBuyAmount = totalBuyAmount
			} else {
				ts, totalSellAmount := topSellParser(s)
				lInfo.TopSell = fmt.Sprintf("[%s]", strings.Join(ts, ","))
				lInfo.TotalSellAmount = totalSellAmount
			}
		})
		lInfos = append(lInfos, lInfo)
	})
	return lInfos
}

// 上榜详细原因, 一支股票可能上榜多次, 所有会有多个原因
func (i *info) reasonHandle(contentBox *goquery.Selection) []string {
	var reasons []string
	contentBox.Find(".content .data-tips").Each(func(i int, s *goquery.Selection) {
		reasons = append(reasons, s.Find(".left").Text())
	})
	return reasons
}

func (i *info) Do() {
	var results []EastMoney
	i.db.Select("s_code", "s_name", "tdate", "ctypedes", "dp", "top_buy", "top_sell").Find(&results)
	wg := &sync.WaitGroup{}
	for _, result := range results {
		tb := result.TopBuy
		ts := result.TopSell
		if tb == "" && ts == "" {
			wg.Add(1)
			go i.fetch(result, wg)
			time.Sleep(time.Millisecond * 240)
		}
	}
	wg.Wait()
}

// topSellParser: parse daily five top sell
func topSellParser(s *goquery.Selection) ([]string, float64) {
	var topSellArr []string
	var totalSellAmount float64

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
		topSellArr = append(topSellArr, string(ts))
		totalSellAmount += topSell.SellAmount
	}
	return topSellArr, totalSellAmount
}

// topBuyParser: parse daily five top buy
func topBuyParser(s *goquery.Selection) ([]string, float64) {
	var topBuyArr []string
	var totalBuyAmount float64
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
		topBuyArr = append(topBuyArr, string(tb))
		totalBuyAmount += topBuy.BuyAmount
	}
	return topBuyArr, totalBuyAmount
}

// findDealInfo: find deal information
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

// findSCName : find stock name and the info link
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
