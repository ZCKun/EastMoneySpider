package utils

import (
	"log"
	"regexp"
	"strings"
)

// infoReason 上榜原因明细
func ReasonContains(listReason, infoReason string) bool {
	var text string

	reasonRegex, _ := regexp.Compile(`^(\d{4}-\d{2}-\d{2}).*?类型：(.*?)$`)
	turRegex, _ := regexp.Compile(`(.*?)的.*证券`)

	results := reasonRegex.FindStringSubmatch(infoReason)
	reason := strings.ReplaceAll(listReason, "达到", "达")

	if blocks := strings.Split(reason, "，"); len(blocks) > 1 {
		m := 0
		for _, block := range blocks {
			if i := len(block); i > m {
				text = block
				m = i
			}
		}
		text = strings.Split(text, "达")[0]
	} else if strings.Contains(reason, "换手率") {
		seq := "换手率"
		text = seq + strings.Split(reason, seq)[1]
		match := turRegex.FindStringSubmatch(text)
		if len(match) <= 1 {
			log.Printf("ReasonContains:error:30line:cant find seq from text. listReason: %s, infoReason: %s\n",
				listReason, infoReason)
			return false
		}
		// throw exception
		text = match[1]
	} else {
		if strings.Contains(reason, "跌幅") {
			seq := "跌幅"
			text = seq + strings.Split(reason, seq)[1]
		} else if strings.Contains(reason, "涨幅") {
			seq := "涨幅"
			text = seq + strings.Split(reason, seq)[1]
		}
	}

	if strings.Contains(text, "的证券") {
		text = strings.ReplaceAll(text, "的证券", "")
	}

	return strings.Contains(results[2], text)
}
