package search

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Utils struct{}


func (u *Utils) request(req *http.Request, retry int) (string, error) {
	client := &http.Client{}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	url := req.URL.String()

	if resp.StatusCode == 502 || resp.StatusCode == 504 {
		// 出现502或504, 尝试5次重新请求
		if retry > 5 {
			return "", fmt.Errorf("Retry ListRequest Url %s HTTP Response Error %d 5 times.\n", url, resp.StatusCode)
		}
		time.Sleep(time.Second)
		return u.request(req, retry+1)
	}

	// respBody := transform.NewReader(resp.Body, simplifiedchinese.GB18030.NewDecoder())
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func (u *Utils) ListRequest(params map[string]string) (string, error) {
	url := "http://datainterface3.eastmoney.com//EM_DataCenter_V3/api/LHBGGDRTJ/GetLHBGGDRTJ"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Host", "datainterface3.eastmoney.com")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	return u.request(req, 1)
}

func (u *Utils) InfoRequest(date string, symbolCode string) (string, error) {
	url := fmt.Sprintf("https://data.eastmoney.com/stock/lhb,%s,%s.html", date, symbolCode)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Referer", "https://data.eastmoney.com/stock/tradedetail.html")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Host", "data.eastmoney.com")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	return u.request(req, 1)
}

// 获取请求页面url. startDate 开始日期, endDate 结束日期, page 页
func (u *Utils) GetUrl(startDate, endDate string, page int) string {
	var url bytes.Buffer
	url.WriteString("http://data.eastmoney.com/DataCenter_V3/stock2016/TradeDetail/")
	params := fmt.Sprintf(
		"pagesize=%d,page=%d,sortRule=-1,sortType=,startDate=%s,endDate=%s,gpfw=0,js=var%%20data_tab_2.html",
		50, page, startDate, endDate)
	url.WriteString(params)
	return url.String()
}

func (u *Utils) GetParams(startDate, endDate string, page int) map[string]string {
	return map[string]string{
		"tkn":           "eastmoney",
		"mkt":           "0",
		"dateNum":       "",
		"startDateTime": startDate,
		"endDateTime":   endDate,
		"sortRule":      "1",
		"sortColumn":    "",
		"pageNum":       strconv.Itoa(page),
		"pageSize":      "50",
		"cfg":           "lhbggdrtj",
	}
}