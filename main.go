package main

import (
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/gomail.v2"
	"log"
	"regexp"
)
const (
	FundRootUrl = "http://fundgz.1234567.com.cn/js/"
)


func FetchFund(codes []string) (fund []map[string]string) {
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
			fund = append(fund, fundDataMap)
			table := `
            <tr>
              <td width="50" align="center">` + fundDataMap["name"] + `</td>
              <td width="50" align="center">` + fundDataMap["gszzl"] + `%</td>
              <td width="50" align="center">` + fundDataMap["gsz"] + `</td>
              <td width="50" align="center">` + fundDataMap["dwjz"] + `</td>
              <td width="50" align="center">` + fundDataMap["gztime"] + `</td>
            </tr>
		`
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
							</tr>` + table + `
						</table>
					</div>
            	</div>
            </div>
            </body>
        </html>`
			return html
		}
	}
	return ""
}

func GenHTML(funds []map[string]string) {
	for key, value := range funds {

	}
}

func SendEmail(fundData string)  {
	if fundData == "" {
		return
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "365999802@qq.com")
	m.SetHeader("To", "365999802@qq.com")
	m.SetHeader("Subject", "基金涨跌监控")
	m.SetBody("text/html", fundData)

	d := gomail.NewDialer("smtp.qq.com", 587, "365999802@qq.com", "")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func main() {
	//ret := FetchFund("180012")
	fundSlice := []string{"180012", "167301"}
	ret := FetchFund(fundSlice)
	SendEmail(ret)
}
