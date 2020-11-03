package utils

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type HttpUtil struct{}

func (u *HttpUtil) request(req *http.Request, retry int) (string, error) {
	client := &http.Client{}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			if retry < 10 {
				time.Sleep(time.Second)
				return u.request(req, retry+1)
			}
		}
		return "", err
	}

	defer resp.Body.Close()

	url := req.URL.String()

	if resp.StatusCode == 502 || resp.StatusCode == 504 {
		// 出现502或504, 尝试5次重新请求
		if retry > 20 {
			return "", fmt.Errorf("Retry ListRequest Url %s HTTP Response Error %d 5 times.\n", url, resp.StatusCode)
		}
		time.Sleep(time.Second)
		return u.request(req, retry+1)
	}

	body := resp.Body
	content, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// 榜单
func (u *HttpUtil) ListRequest(params map[string]string) (string, error) {
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
	body, err := u.request(req, 1)
	if err != nil {
		return req.URL.RawQuery, err
	}
	return body, nil
}

// 明细
func (u *HttpUtil) InfoRequest(date string, symbolCode string) (string, error) {
	url := fmt.Sprintf("https://data.eastmoney.com/stock/lhb,%s,%s.html", date, symbolCode)
	fmt.Println(url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Referer", "https://data.eastmoney.com/stock/tradedetail.html")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Host", "data.eastmoney.com")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	body, err := u.request(req, 1)
	c, _ := ioutil.ReadAll(transform.NewReader(strings.NewReader(body), simplifiedchinese.GB18030.NewDecoder()))
	body = string(c)
	if err != nil {
		return url, err
	}
	return body, nil
}

// 获取请求页面url. startDate 开始日期, endDate 结束日期, page 页
func (u *HttpUtil) GetUrl(startDate, endDate string, page int) string {
	var url bytes.Buffer
	url.WriteString("http://data.eastmoney.com/DataCenter_V3/stock2016/TradeDetail/")
	params := fmt.Sprintf(
		"pagesize=%d,page=%d,sortRule=-1,sortType=,startDate=%s,endDate=%s,gpfw=0,js=var%%20data_tab_2.html",
		50, page, startDate, endDate)
	url.WriteString(params)
	return url.String()
}

func (u *HttpUtil) GetParams(startDate, endDate string, page int) map[string]string {
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
