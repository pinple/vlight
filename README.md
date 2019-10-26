# vlight![](https://travis-ci.com/Neulana/vlight.svg?token=ATFZGq5q9tbJu6KjZsyy&branch=master)
基金涨跌监控，每天14:30定时抓取基金行情，如果超过自己设置的阈值，则发邮件给自己，根据行情做相应的加仓/建仓动作。解放你的注意力，避免因为频繁查看行情而**影响自己的情绪**，从而冲动操作，买基金最重要就是：选好一支基金并且长期持有。
# 如何构建
## 一、无服务器
使用[Travis](https://www.travis-ci.org/)构建，在Travis后台设置邮箱的环境变量，EMAIL_NAME和EMAIL_PASSWORD，然后将构建时间设置为【每天】。
## 二、有自己的服务器
建cron任务。