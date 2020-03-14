package main

import (
	//"encoding/json"
	//"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	//"io/ioutil"
	//"strconv"
)


func ExitHandler(whs *WeedHttpServer){
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		log.Println("Ctrl+C pressed in Terminal, Exit Slaves and Server Monitor...")
		for idx := range (*whs).Slaves{
			whs.Slaves[idx].Exit()
		}
		if whs.MonitorServer{
			whs.TargetServer.ExitTargetServer()
		}
		os.Exit(0)
	}()
}

func main() {
	var whs WeedHttpServer
	whs.Distribute = true
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/get_second_data", whs.getSecondData)
	http.HandleFunc("/slave_report", whs.SlaveReport)
	http.HandleFunc("/slave_report_exit", whs.SlaveReportExit)
	http.HandleFunc("/target_server_report", whs.TargetServerReport)
	http.HandleFunc("/target_server_report_exit", whs.TargetServerReportExit)
	http.HandleFunc("/start", whs.StartTestHandler)
	http.HandleFunc("/stop", whs.StopTestHandler)
	http.HandleFunc("/get_server_second_data", whs.GetTargetServerData)
	http.Handle("/static/", http.FileServer(http.Dir("./")))
	ExitHandler(&whs)
	http.ListenAndServe(":6767", nil)
}
