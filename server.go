package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

func ExitHandler(master_ip string, json_str string) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		log.Println("Ctrl+C pressed in Terminal, exit server monitor...")
		req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s:6767/target_server_report_exit", master_ip), strings.NewReader(json_str))
		client := &http.Client{}
		client.Do(req)
		os.Exit(0)
	}()
}

func main() {
	check_port_conn, _ := net.Listen("tcp", ":0")
	port := check_port_conn.Addr().(*net.TCPAddr).Port
	check_port_conn.Close()
	//log.Printf("%d",port)
	var master_ip string
	flag.StringVar(&master_ip, "master-ip", "127.0.0.1", "specify the master ip address")
	flag.Parse()
	addrs, _ := net.InterfaceAddrs()
	var ip_addr string
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && !strings.HasPrefix(ipnet.IP.String(), "169") {
				ip_addr = ipnet.IP.String()
				break
			}
		}
	}
	now := int64(time.Now().UnixNano() / 1e6)
	var target_server TargetServer
	target_server.IPAddr = fmt.Sprintf("%s:%d", ip_addr, port)
	target_server.TimeStampInMs = now
	json_str, err := json.Marshal(target_server)
	if err != nil {
		fmt.Errorf("Marshal Error %v", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:6767/target_server_report", master_ip), strings.NewReader(string(json_str)))
	if err != nil {
		log.Println(err)
		return
	}
	client := &http.Client{}
	var resp *http.Response
	var err_http error
	for {
		resp, err_http = client.Do(req)
		if err_http != nil {
			log.Println("Try connecting master error: ", err_http)
			time.Sleep(time.Second)
			continue
		}
		break
	}
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
	resp.Body.Close()
	var whs WeedHttpServer
	http.HandleFunc("/get_server_second_data", whs.GetTargetServerData)
	http.HandleFunc("/exit", whs.ExitHandler)
	ExitHandler(master_ip, string(json_str))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
