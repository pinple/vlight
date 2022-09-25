package main

const (
	FundJsUrl    = "http://fundgz.1234567.com.cn/js/"
	FundHTMLUrl  = "http://fund.eastmoney.com/"
	MIN_RISE_NUM = 0.1
	MAX_FALL_NUM = -0.1

	DailyTitle = `
                 <tr>
	             <td width="50" align="center">基金名称</td>
	             <td width="50" align="center">估算涨幅</td>
	             <td width="50" align="center">当前估算净值</td>
	             <td width="50" align="center">昨日单位净值</td>
	             <td width="50" align="center">估算时间</td>
                 </tr>
                 `
	WeeklyTitle = `
                 <tr>
	             <td width="50" align="center">基金名称</td>
	             <td width="50" align="center">近1周净值变化</td>
                 </tr>
                 `
	MonthlyTitle = `
                 <tr>
	             <td width="50" align="center">基金名称</td>
	             <td width="50" align="center">近1月净值变化</td>
                 </tr>
                 `
)
