package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"time"
)

type Task struct {
	/*任务结构体，因部分测试场景需要共享数据，故使用此结构体避免全局变量的使用，如本demo的场景需要使用长连接测试
	*/
	client	http.Client
}

func (task *Task)Init() {
	/*此函数，每次开始测试只会执行一次(Init执行一次、Work执行多次)，用于初始化
	*/
	task.client = http.Client{ Timeout : time.Second }
}

func (task *Task)Work() (bool, int64) {
	/*此函数为压测具体事项，如下为GET腾讯首页的示例
	*/
	req, _ := http.NewRequest("GET", "https://www.qq.com", strings.NewReader(""))
	var start_time, cost_time int64
	start_time = int64(time.Now().UnixNano() / 1e6)
	resp, err := task.client.Do(req) //此处使用了上面初始化的client
	ioutil.ReadAll(resp.Body)
	end_time := int64(time.Now().UnixNano() / 1e6)
	cost_time = end_time - start_time
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		fmt.Println(err)
		return false, 0
	}
	if resp.StatusCode != 200 {
		fmt.Println(resp.Header)
		fmt.Printf("Expect 200, %d GET!", resp.StatusCode)
		return false, 0
	}
	return true, int64(cost_time)
	//fmt.Println(resp)
}

func (task *Task)Uninit() {
	/*此函数，每一次压测只会执行一次(同Init)，用于释放初始化的资源，这个示例没有资源需要释放
	*/
}

