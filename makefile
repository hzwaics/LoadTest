.PHONY: weeds
weeds:
	go build -o weeds work.go core.go http_server.go single.go
all: slave master server-monitor weeds
slave:
	go build -o weeds-slave http_server.go slave.go core.go work.go
master:
	go build -o weeds-master core.go http_server.go master.go work.go
server-monitor:
	go build -o weeds-server-monitor work.go core.go http_server.go server.go
clean:
	rm -f weeds weeds-slave weeds-master
