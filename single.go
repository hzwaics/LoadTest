package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
)

func ExitHandler(whs *WeedHttpServer){
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		log.Println("Ctrl+C pressed in Terminal, exit server monitor...")
		if whs.MonitorServer{
			whs.TargetServer.ExitTargetServer()
		}
		os.Exit(0)
	}()
}
func main() {
	var whs WeedHttpServer
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/get_second_data", whs.getSecondData)
	http.HandleFunc("/slave_report", whs.SlaveReport)
	http.HandleFunc("/target_server_report", whs.TargetServerReport)
	http.HandleFunc("/target_server_report_exit", whs.TargetServerReportExit)
	http.HandleFunc("/start", whs.StartTestHandler)
	http.HandleFunc("/stop", whs.StopTestHandler)
	http.Handle("/static/", http.FileServer(http.Dir("./")))
	ExitHandler(&whs)
	http.ListenAndServe(":6767", nil)
}
