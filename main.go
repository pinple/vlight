package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/gomail.v2"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	fundCodeStr, watchWeekDay, watchMonthDay, emailName, emailPassword string
)
var (
	fundCodeSlice []string
)

func init() {
	fundCodeStr = os.Getenv("WATCH_FUND_CODE")
	//fundCodeStr = "005827,012414,003095,161005"
	fundCodeSlice = strings.Split(fundCodeStr, ",")
	watchWeekDay = os.Getenv("WATCH_WEEK_DAY")
	watchMonthDay = os.Getenv("WATCH_MONTH_DAY")
	emailName = os.Getenv("EMAIL_NAME")
	emailPassword = os.Getenv("EMAIL_PASSWORD")
}

func fetchFund(codes []string) []map[string]string {
	var fundResult []map[string]string
	var weeklyChange string
	var oneMonthChange string
	fmt.Printf("codes: %#v\n", codes)
	for _, code := range codes {
		fundJsUrl := FundJsUrl + code + ".js"
		resp1, err := http.Get(fundJsUrl)
		if err != nil {
			panic(err)
		}
		defer resp1.Body.Close()

		fundHTMLUrl := FundHTMLUrl + code + ".html"
		fmt.Printf("fundHTMLUrl: %s\n", string(fundHTMLUrl))
		resp2, err := http.Get(fundHTMLUrl)
		if err != nil {
			panic(err)
		}
		defer resp2.Body.Close()

		re, _ := regexp.Compile("jsonpgz\\((.*)\\);")
		data, err := io.ReadAll(resp1.Body)
		if err != nil {
			panic(err)
		}
		ret := re.FindSubmatch(data)
		fundData := ret[1]
		fmt.Printf("fundData: %s\n", string(fundData))

		doc, err2 := goquery.NewDocumentFromReader(resp2.Body)
		if err != nil {
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

	fmt.Printf("fundResult: %#v\n", fundResult)
	return fundResult
}

func sendEmail(content string) {
	if content == "" {
		return
	}
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
	fundResult := fetchFund(fundCodeSlice)
	content := renderHTML(fundResult)
	sendEmail(content)
}
