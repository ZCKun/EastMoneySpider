package search

import (
	"gorm.io/gorm"
)

type (
	LHBList struct {
		Bmoney      string `json:"Bmoney"`
		Chgradio    string `json:"Chgradio"`
		ClosePrice  string `json:"ClosePirce"`
		Ctypedes    string `json:"Ctypedes"`
		DP          string `json:"DP"`
		Dchratio    string `json:"Dchratio"`
		JD          string `json:"JD"`
		JGBMoney    string `json:"JGBMoney"`
		JGBSumCount string `json:"JGBSumCount"`
		JGJMMoney   string `json:"JGJMMoney"`
		JGSMoney    string `json:"JGSMoney"`
		JGSSumCount string `json:"JGSSumCount"`
		JmMoney     string `json:"JmMoney"`
		JmRate      string `json:"JmRate"`
		Ltsz        string `json:"Ltsz"`
		Ntransac    string `json:"Ntransac"`
		Oldid       string `json:"Oldid"`
		Rchange1dc  string `json:"Rchange1dc"`
		Rchange1do  string `json:"Rchange1do"`
		Rchange1m   string `json:"Rchange1m"`
		Rchange1y   string `json:"Rchange1y"`
		Rchange2dc  string `json:"Rchange2dc"`
		Rchange2do  string `json:"Rchange2do"`
		Rchange3dc  string `json:"Rchange3dc"`
		Rchange3do  string `json:"Rchange3do"`
		Rchange3m   string `json:"Rchange3m"`
		Rchange5dc  string `json:"Rchange5dc"`
		Rchange5do  string `json:"Rchange5do"`
		Rchange6m   string `json:"Rchange6m"`
		Rchange10dc string `json:"Rchange10dc"`
		Rchange10do string `json:"Rchange10do"`
		Rchange15dc string `json:"Rchange15dc"`
		Rchange15do string `json:"Rchange15do"`
		Rchange20dc string `json:"Rchange20dc"`
		Rchange20do string `json:"Rchange20do"`
		Rchange30dc string `json:"Rchange30dc"`
		Rchange30do string `json:"Rchange30do"`
		SCode       string `json:"SCode"`
		SName       string `json:"SName"`
		Smoney      string `json:"Smoney"`
		SumCount    string `json:"SumCount"`
		Tdate       string `json:"Tdate"`
		Turnover    string `json:"Turnover"`
		ZeMoney     string `json:"ZeMoney"`
		ZeRate      string `json:"ZeRate"`
	}

	DealInfo struct {
		SCName                  string  `json:"sc_name"`
		Href                    string  `json:"href"`
		BuyAmount               float64 `json:"buy_amount"`
		BuyAmountPropTotalTran  string  `json:"buy_amount_prop_total_tran"`
		SellAmount              float64 `json:"sell_amount"`
		SellAmountPropTotalTran string  `json:"sell_amount_prop_total_tran"`
		NetAmount               float64 `json:"net_amount"`
	}

	LHBInfo struct {
		TopSell         string
		TopBuy          string
		TotalBuyAmount  float64
		TotalSellAmount float64
	}

	EastMoney struct {
		*LHBList
		*LHBInfo
	}
)

func Search(db *gorm.DB) {
	i := &info{db}
	i.Do()
}
