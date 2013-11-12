package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type ServiceType int

const (
	FILESYSTEM ServiceType = iota
	DIRECTORY
	FILE
	PROCESS
	HOST
	SYSTEM
	FIFO
	PROGRAM
)

type Monit struct {
	Id          string    `xml:"id,attr"          json:"id"`
	Incarnation int       `xml:"incarnation,attr" json:"incarnation"`
	Version     string    `xml:"version,attr"     json:"version"`
	IpAddress   string    `                       json:"ip_address"`
	Server      Server    `xml:"server"           json:"server"`
	Platform    Platform  `xml:"platform"         json:"platform"`
	Service     []Service `xml:"services>service" json:"services"`
}

type Server struct {
	Uptime        int    `xml:"uptime"               json:"uptime"`
	Poll          int    `xml:"poll"                 json:"poll"`
	StartDelay    int    `xml:"startdelay"           json:"start_delay"`
	Localhostname string `xml:"localhostname"        json:"localhostname"`
	HttpAddress   string `xml:"httpd>address"        json:"-"`
	HttpPort      int    `xml:"httpd>port"           json:"-"`
	Username      string `xml:"credentials>username" json:"-"`
	Password      string `xml:"credentials>password" json:"-"`
}

type Platform struct {
	Name    string `xml:"name"    json:"name"`
	Release string `xml:"release" json:"release"`
	Version string `xml:"version" json:"version"`
	Machine string `xml:"machine" json:"machine"`
	Cpu     int    `xml:"cpu"     json:"cpu"`
	Memory  int    `xml:"memory"  json:"memory"`
	Swap    int    `xml:"swap"    json:"swap"`
}

type Service struct {
	Name          string      `xml:"name,attr"      json:"name"`
	Type          ServiceType `xml:"type"           json:"type"`
	CollectedSec  int64       `xml:"collected_sec"  json:"-"`
	CollectedUsec int64       `xml:"collected_usec" json:"-"`
	Status        int         `xml:"status"         json:"status"`
	Monitor       int         `xml:"monitor"        json:"monitor"`
	MonitorMode   int         `xml:"monitormode"    json:"monitor_mode"`
	PendingAction int         `xml:"pendingaction"  json:"pending_action"`
	StatusMessage string      `xml:"status_message" json:"status_message,omitempty"`
}

func MonitFromXml(r io.Reader) Monit {
	var monit Monit
	d := xml.NewDecoder(r)
	d.CharsetReader = CharsetReader
	d.Decode(&monit)
	return monit
}

func (server *Server) Execute(service string, action string) {
	params := url.Values{}
	params.Set("action", action)
	url := fmt.Sprintf("http://%s:%d/%s", server.HttpAddress, server.HttpPort, service)
	req, _ := http.NewRequest("POST", url, strings.NewReader(params.Encode()))
	req.SetBasicAuth(server.Username, server.Password)

	_, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Server[%s]: Error: %s", server.Localhostname, err)
	}
}
