package main

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type OneSecondData struct {
	Req               int64 `json:"request_num"`
	Succ_resp         int64 `json:"success_response_num"`
	Average_cost_time int64 `json:"average_cost_time"`
	Failed_num        int64 `json:"failed_num"`
	Time_stamp        int64 `json:"time_stamp"`
}

type Weeds struct {
	request_num_in_one_second       int64
	succ_response_num_in_one_second int64
	total_request_num               int64
	total_succ_response_num         int64
	total_failed_num                int64
	total_resp_time_in_one_second   int64
	failed_num_in_one_second        int64
	sleep_time_in_microsecond       int64
	Status                          int64 //0:stop, 1:start
	last_second_data                OneSecondData
	task				Task
}

func (weeds *Weeds) StopLoadTest() {
	weeds.Status = 0
	atomic.StoreInt64(&(weeds.request_num_in_one_second), 0)
	atomic.StoreInt64(&(weeds.succ_response_num_in_one_second), 0)
	atomic.StoreInt64(&(weeds.total_resp_time_in_one_second), 0)
	atomic.StoreInt64(&(weeds.failed_num_in_one_second), 0)
	weeds.task.Uninit()
}

func (weeds *Weeds) RunSingleJob() {
	succ, cost_time := weeds.task.Work()
	if succ {
		atomic.AddInt64(&(weeds.succ_response_num_in_one_second), 1)
		atomic.AddInt64(&(weeds.total_succ_response_num), 1)
		atomic.AddInt64(&(weeds.total_resp_time_in_one_second), int64(cost_time))
	} else {
		atomic.AddInt64(&(weeds.failed_num_in_one_second), 1)
		atomic.AddInt64(&(weeds.total_failed_num), 1)
	}
}

func (weeds *Weeds) RunJobs(qps float64, wg *sync.WaitGroup) {
	for weeds.Status == 1 {
		time.Sleep(time.Duration(weeds.sleep_time_in_microsecond) * time.Microsecond)
		atomic.AddInt64(&(weeds.request_num_in_one_second), 1)
		atomic.AddInt64(&(weeds.total_request_num), 1)
		go weeds.RunSingleJob()
	}
	wg.Done()
}

func (weeds *Weeds) DelayTimeAdjust(qps float64) {
	time.Sleep(time.Second)
	for weeds.Status == 1 && weeds.last_second_data.Req < int64(qps) && weeds.sleep_time_in_microsecond > 5 {
		time.Sleep(time.Duration(1000) * time.Millisecond)
		weeds.sleep_time_in_microsecond -= int64(0.05 * float64(weeds.sleep_time_in_microsecond))
	}
	for weeds.Status == 1 {
		time.Sleep(time.Second)
		if weeds.last_second_data.Req > int64(qps) {
			weeds.sleep_time_in_microsecond += int64(0.01 * float64(weeds.sleep_time_in_microsecond))
		} else if weeds.last_second_data.Req < int64(qps) {
			weeds.sleep_time_in_microsecond -= int64(0.01 * float64(weeds.sleep_time_in_microsecond))
		}
	}
}

func (weeds *Weeds) Controller() {
	for weeds.Status == 1 {
		time.Sleep(time.Second)
		var average_resp_time int64
		if atomic.LoadInt64(&(weeds.succ_response_num_in_one_second)) == 0 {
			average_resp_time = 0
		} else {
			average_resp_time = int64(atomic.LoadInt64(&(weeds.total_resp_time_in_one_second)) / int64(atomic.LoadInt64(&(weeds.succ_response_num_in_one_second))))
		}
		//fmt.Printf("total req:%d, total resp:%d, req/s:%d, resp/s:%d, average_resp_time:%d\n", weeds.total_request_num, weeds.total_succ_response_num, weeds.request_num_in_one_second, weeds.succ_response_num_in_one_second, average_resp_time)
		now := int64(time.Now().Unix())
		weeds.last_second_data = OneSecondData{weeds.request_num_in_one_second, weeds.succ_response_num_in_one_second, average_resp_time, weeds.failed_num_in_one_second, now}
		atomic.StoreInt64(&(weeds.request_num_in_one_second), 0)
		atomic.StoreInt64(&(weeds.succ_response_num_in_one_second), 0)
		atomic.StoreInt64(&(weeds.total_resp_time_in_one_second), 0)
		atomic.StoreInt64(&(weeds.failed_num_in_one_second), 0)
	}
}

func (weeds *Weeds) GetTestDataSum() TestDataSum {
	var test_data_sum TestDataSum
	test_data_sum.TotalReqNum = weeds.total_request_num
	test_data_sum.TotalRespNum = weeds.total_succ_response_num
	test_data_sum.TotalFailedNum = weeds.total_failed_num
	return test_data_sum
}

func (weeds *Weeds) GetSecondData() OneSecondData {
	if weeds.Status == 1 {
		return weeds.last_second_data
	} else {
		var zero OneSecondData
		zero.Time_stamp = int64(time.Now().Unix())
		return zero
	}
}

func (weeds *Weeds) StopInTime(duration_time int64) {
	time.Sleep(time.Duration(duration_time) * time.Second)
	weeds.Status = 0
}

func (weeds *Weeds) QpsController(start_qps float64, end_qps float64, qps_step float64, qps *float64, cpu_num int) {
	for *qps < end_qps {
		time.Sleep(time.Duration(1) * time.Second)
		weeds.sleep_time_in_microsecond = int64(float64(1000000) * float64(cpu_num) / *qps)
		*qps = *qps + qps_step
	}
	weeds.DelayTimeAdjust(*qps)
}

func (weeds *Weeds) StartLoadTest(start_qps float64, end_qps float64, qps_step float64, duration_time int64) {
	weeds.Status = 1
	if duration_time > 0 {
		go weeds.StopInTime(duration_time)
	}
	//go weeds.DelayTimeAdjust(qps)
	if qps_step <= 0 {
		qps_step = end_qps
	}
	qps := start_qps
	if qps == 0 {
		qps = qps_step
	}
	cpu_num := runtime.NumCPU()
	var wg sync.WaitGroup
	if int64(qps) < int64(cpu_num*10) {
		cpu_num = 1 //qps太小的话，当作单核处理
	}
	go weeds.QpsController(start_qps, end_qps, qps_step, &qps, cpu_num)
	go weeds.Controller()
	weeds.sleep_time_in_microsecond = int64(float64(1000000) * float64(cpu_num) / qps)
	weeds.task.Init()
	for t := 0; t < cpu_num; t++ {
		wg.Add(1)
		go weeds.RunJobs(qps, &wg)
	}
	wg.Wait()
	return
}
