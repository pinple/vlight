package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	FundJsUrl       = "http://fundgz.1234567.com.cn/js/"
	FundHTMLUrl     = "http://fund.eastmoney.com/"
	MIN_RISE_NUM    = 1.5
	MAX_FALL_NUM    = -1.5
)

var fundCodeStr = os.Getenv("WATCH_FUND_CODE")
var fundCodeSlice = strings.Split(fundCodeStr, "|")

var dailyTitle = `
                 <tr>
	             <td width="50" align="center">基金名称</td>
	             <td width="50" align="center">估算涨幅</td>
	             <td width="50" align="center">当前估算净值</td>
	             <td width="50" align="center">昨日单位净值</td>
	             <td width="50" align="center">估算时间</td>
                 </tr>
                 `

var weeklyTitle = `
                 <tr>
	             <td width="50" align="center">基金名称</td>
	             <td width="50" align="center">近1周净值变化</td>
                 </tr>
                 `

var oneMonthTitle = `
                 <tr>
	             <td width="50" align="center">基金名称</td>
	             <td width="50" align="center">近1月净值变化</td>
                 </tr>
                 `

func FetchFund(codes []string) []map[string]string {
	var fundResult []map[string]string
	var weeklyChange string
	var oneMonthChange string
	for _, code := range codes {
		fundJsUrl := FundJsUrl + code + ".js"
		request := gorequest.New()
		resp, body, err := request.Get(fundJsUrl).End()
		defer resp.Body.Close()
		if err != nil {
			log.Fatal(err)
			return nil
		}

		fundHTMLUrl := FundHTMLUrl + code + ".html"
		resp1, body1, err1 := request.Get(fundHTMLUrl).End()
		defer resp1.Body.Close()
		if err1 != nil {
			log.Fatal(err1)
			return nil
		}

		re, _ := regexp.Compile("jsonpgz\\((.*)\\);")
		ret := re.FindSubmatch([]byte(body))
		fundData := ret[1]

		doc, err2 := goquery.NewDocumentFromReader(strings.NewReader(body1))
		if err2 != nil {
			log.Fatal(err2)
		}
		doc.Find("#increaseAmount_stage > table:nth-child(1) > tbody:nth-child(1) > tr:nth-child(2)").Each(func(i int, s *goquery.Selection) {
			s.Find("td > div").Each(func(j int, k *goquery.Selection) {
				change := k.Text()
				switch j {
				case 1:
					weeklyChange = change
				case 2:
					oneMonthChange = change
				}
			})
		})
		var fundDataMap map[string]string
		if err := json.Unmarshal(fundData, &fundDataMap); err == nil {
			fundDataMap["weeklyChange"] = weeklyChange
			fundDataMap["oneMonthChange"] = oneMonthChange
			fundResult = append(fundResult, fundDataMap)
		}
	}
	return fundResult
}

func GenerateHTML(fundResult []map[string]string) string {
	var watchWeekDay = os.Getenv("WATCH_WEEK_DAY")
	var watchMonthDay = os.Getenv("WATCH_MONTH_DAY")
	var dailyElements []string
	var dailyContent string
	var weeklyElements []string
	var weeklyContent string
	var oneMonthElements []string
	var oneMonthContent string
	var dailyText string
	var weeklyText string
	var oneMonthText string
	now := time.Now()
	for _, fund := range fundResult {
		gszzl, err := strconv.ParseFloat(fund["gszzl"], 32)
		if err != nil {
			fmt.Printf("!!error!!: %s", err)
			continue
		}
		if gszzl > 0 {
			fund["gszzl"] = "+" + strconv.FormatFloat(gszzl, 'f', -1, 32)
		}
		// 每日涨幅，涨跌幅度超出设定值才发出通知
		if (gszzl > 0 && gszzl >= MIN_RISE_NUM) || (gszzl < 0 && gszzl <= MAX_FALL_NUM) {
			dailyElement := `
                                   <tr>
                                     <td width="50" align="center">` + fund["name"] + `</td>
                                     <td width="50" align="center">` + fund["gszzl"] + `%</td>
                                     <td width="50" align="center">` + fund["gsz"] + `</td>
                                     <td width="50" align="center">` + fund["dwjz"] + `</td>
                                     <td width="50" align="center">` + fund["gztime"] + `</td>
                                   </tr>
	                           `
			dailyElements = append(dailyElements, dailyElement)
		}
		// 一周涨幅
		if now.Weekday().String() == watchWeekDay {
			weeklyElement := `
                                   <tr>
                                     <td width="50" align="center">` + fund["name"] + `</td>
                                     <td width="50" align="center">` + fund["weeklyChange"] + `</td>
                                   </tr>
                                   `
			weeklyElements = append(weeklyElements, weeklyElement)
		}
		// 月度涨幅
		monthNum, err := strconv.Atoi(watchMonthDay)
		if now.Day() == monthNum {
			oneMonthElement := `
                                   <tr>
                                     <td width="50" align="center">` + fund["name"] + `</td>
                                     <td width="50" align="center">` + fund["oneMonthChange"] + `</td>
                                   </tr>
                                   `
			oneMonthElements = append(oneMonthElements, oneMonthElement)
		}
	}
	dailyContent = strings.Join(dailyElements, "\n")
	weeklyContent = strings.Join(weeklyElements, "\n")
	oneMonthContent = strings.Join(oneMonthElements, "\n")
	if dailyContent != "" || weeklyContent != "" || oneMonthContent != "" {
		if dailyContent != "" {
			dailyText = `
                                    <table width="30%" border="1" cellspacing="0" cellpadding="0">
				    ` + dailyTitle + dailyContent + `
				    </table> <br><br>`
		}
		if weeklyContent != "" {
			weeklyText = `
                                    <table width="30%" border="1" cellspacing="0" cellpadding="0">
				    ` + weeklyTitle + weeklyContent + `
				    </table> <br><br>`
		}
		if oneMonthContent != "" {
			oneMonthText = `
                                    <table width="30%" border="1" cellspacing="0" cellpadding="0">
				    ` + oneMonthTitle + oneMonthContent + `
				    </table> <br><br>`
		}
		html := `
			</html>
			    <head>
			        <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
			    </head>
                            <body>
			        <div id="container">
			            <p>基金涨跌监控:</p>
			            <div id="content">
				            ` + dailyText + weeklyText + oneMonthText + `
				    </div>
            	                </div>
                            </body>
                        </html>`

		return html
	}

	return ""
}

func SendEmail(content string) {
	if content == "" {
		return
	}
	emailName := os.Getenv("EMAIL_NAME")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	m := gomail.NewMessage()
	m.SetHeader("From", emailName)
	m.SetHeader("To", emailName)
	m.SetHeader("Subject", "基金涨跌监控")
	m.SetBody("text/html", content)
	d := gomail.NewDialer("smtp.qq.com", 587, emailName, emailPassword)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func main() {
	fundResult := FetchFund(fundCodeSlice)
	content := GenerateHTML(fundResult)
	SendEmail(content)
}
