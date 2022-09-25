package main

import (
	"strconv"
	"strings"
	"time"
)

func renderHTML(fundResult []map[string]string) string {
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
			panic(err)
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
				    ` + DailyTitle + dailyContent + `
				    </table> <br><br>`
		}
		if weeklyContent != "" {
			weeklyText = `
                                    <table width="30%" border="1" cellspacing="0" cellpadding="0">
				    ` + WeeklyTitle + weeklyContent + `
				    </table> <br><br>`
		}
		if oneMonthContent != "" {
			oneMonthText = `
                                    <table width="30%" border="1" cellspacing="0" cellpadding="0">
				    ` + MonthlyTitle + oneMonthContent + `
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
