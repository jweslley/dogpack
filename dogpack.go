package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

var (
	queue = make(chan *Monit)
	nodes = make(map[string]*Monit)
	port  = flag.Int("p", 2813, "Port to serve dogpack")
)

func bootstrap() {
	for {
		monit := <-queue
		nodes[monit.Server.Localhostname] = monit
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	render(w, "index", tmplData{
		"Title": "Dogpack",
		"Nodes": nodes,
	})
}

func action(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	if req.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	server := req.PostFormValue("server")
	service := req.PostFormValue("service")
	action := req.PostFormValue("action")

	monit, found := nodes[server]
	if found {
		log.Printf("Server[%s]: Executing action '%s' on service '%s'\n", server, action, service)
		go monit.Server.Execute(service, action)
	} else {
		log.Printf("Server[%s]: Unable to execute action '%s' on service '%s'. Server not found.\n", server, action, service)
	}
}

func status(w http.ResponseWriter, req *http.Request) {
	json, err := json.Marshal(nodes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.Write(json)
	}
}

func collector(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	monit := MonitFromXml(req.Body)
	monit.IpAddress, _, _ = net.SplitHostPort(req.RemoteAddr)
	log.Printf("Received status of %s (%s)\n", monit.Server.Localhostname, monit.IpAddress)
	queue <- &monit
}

func main() {
	flag.Parse()

	go bootstrap()

	http.HandleFunc("/", index)
	http.HandleFunc("/status", status)
	http.HandleFunc("/action", action)
	http.HandleFunc("/collector", collector)

	log.Printf("Starting dogpack HTTP server at :%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
