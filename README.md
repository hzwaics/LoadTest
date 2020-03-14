# Weeds
一个golang的压测框架
## 用户指引
(1) 下载工程源码.
(2) 修改work.go, 实现里面的Work()函数, 函数内执行你要的操作, 并返回(bool,int64)类型的两个数值, 第一个数值标记成功(true)与失败(false), 第二个数值为消耗的时间, 单位为毫秒.
(3) 执行` make`如果你没有make，可以复制makefile里面的编译命令进行编译.
(4) 执行` ./weeds`启动工具，打开控制网页.
(5) 如果你要监控目标服务器, 需要把代码复制到目标服务器, 并执行` make server-monitor`编译, 然后执行` ./weeds-server-monitor -master-ip [启动工具所在的机器IP] `启动对目标服务器的监控.
(6) 访问http://[启动工具所在的机器IP]:6767, 管理你的测试任务.
## 分布式执行指引
(1) 下载工程源码到各执行机(slave)和控制机(master).
(2) 在执行机(slave)修改work.go, 实现里面的Work()函数, 函数内执行你要的操作, 并返回(bool,int64)类型的两个数值, 第一个数值标记成功(true)与失败(false), 第二个数值为消耗的时间, 单位为毫秒.
(3) 在控制机(master)执行 `make master`并使用` ./weeds-master`启动控制机(master)
(4) 在执行机(slave)执行`make slave` 并使用`./weeds-slave -master-ip [控制机的机器IP]` 启动执行机(slave)
(5) 如果你要监控目标服务器, 需要把代码复制到目标服务器, 并执行`make server-monitor`编译, 然后执行`./weeds-server-monitor -master-ip [控制机的机器IP]`启动对目标服务器的监控.
(6) 访问http://[控制机的机器IP]:6767, 管理你的测试任务.

## 注意事项
分布式执行时, 控制机(master)和执行机(slave)网络应能相互访问, 如果你要监控服务器资源, 监控的机器和控制机(master)之间也应能相互访问.

# Weeds
A golang load test frame.
## User Guide
(1) Download the source code.
(2) Implent Work() function in work.go .
(3) Run "make".
(4) Run "./weeds" for website control
## Run Distributionally User Guide
(1) Download the source code to all your master and slave machines(If you want to monitor server resource, download it to your server as well).
(2) Run "make master" in your master machine.
(3) Implent Work() function in work.go, then run "make slave" in your slave machines.
(4) If you want to monitor server resource, run "make server-monitor" in the server machine that you about to test.
(5) Run "./weeds-master" on master machine for website control.
(6) Run "./weeds-slave --master-ip [ip or domain of your master]" on your slave machines.
(7) If you want to monitor server resource, run "./weeds-server-monitor" in the server machine that you about to test.
(8) Open the website to control your test, or switch to master machine to get result for no website mode.
