package main

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)
const (
	FundRootUrl = "http://fundgz.1234567.com.cn/js/"
	MIN_RISE_NUM = 1
	MAX_FALL_NUM = -0.8
)

var fundCodeSlice = []string{"180012", "003095", "519778"}


func FetchFund(codes []string) []map[string]string {
	var  fundResult []map[string]string
	for _, code := range codes {
		fundUrl := FundRootUrl + code + ".js"
		request := gorequest.New()
		resp, body, err := request.Get(fundUrl).End()
		defer resp.Body.Close()
		if err != nil {
			log.Fatal(err)
			return nil
		}
		re, _ := regexp.Compile("jsonpgz\\((.*)\\);")
		ret := re.FindSubmatch([]byte(body))
		fundData := ret[1]
		var fundDataMap map[string]string
		if err := json.Unmarshal(fundData, &fundDataMap); err == nil {
			fundResult = append(fundResult, fundDataMap)
		}
	}
	return fundResult
}

func GenerateHTML(fundResult []map[string]string) string {
	var elements []string
	var content string
	for _, fund := range fundResult{
		gszzl, err := strconv.ParseFloat(fund["gszzl"], 32)
		if err != nil {
			fmt.Printf("!!error!!: %s", err)
			continue
		}
		if gszzl > 0 {
			fund["gszzl"] = "+" + strconv.FormatFloat(gszzl, 'f', -1, 32)
		}
		// 涨跌幅度超出设定值才发出通知
		if (gszzl > 0 && gszzl >= MIN_RISE_NUM) || (gszzl < 0 && gszzl <= MAX_FALL_NUM) {
			element := `
            <tr>
              <td width="50" align="center">` + fund["name"] + `</td>
              <td width="50" align="center">` + fund["gszzl"] + `%</td>
              <td width="50" align="center">` + fund["gsz"] + `</td>
              <td width="50" align="center">` + fund["dwjz"] + `</td>
              <td width="50" align="center">` + fund["gztime"] + `</td>
            </tr>
			`
			elements = append(elements, element)
		}
	}
	content = strings.Join(elements, "\n")
	if content != ""{
		html := `
			</html>
				<head>
					<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
				</head>
            <body>
				<div id="container">
					<p>基金涨跌监控:</p>
					<div id="content">
						<table width="30%" border="1" cellspacing="0" cellpadding="0">
							<tr>
							  <td width="50" align="center">基金名称</td>
							  <td width="50" align="center">估算涨幅</td>
							  <td width="50" align="center">当前估算净值</td>
							  <td width="50" align="center">昨日单位净值</td>
							  <td width="50" align="center">估算时间</td>
							</tr>` + content + `
						</table>
					</div>
            	</div>
            </div>
            </body>
        </html>`

		return html
	}

	return ""
}

func SendEmail(content string)  {
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
