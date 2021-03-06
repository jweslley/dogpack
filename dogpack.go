package main

//go:generate esc -o assets.go -prefix=public public

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	queue = make(chan *Monit)
	nodes = make(map[string]*Monit)
	port  = flag.Int("p", 2813, "Port to serve dogpack")
)

func bootstrap() {
	for {
		monit := <-queue
		nodes[monit.Server.HttpAddress] = monit
	}
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
	log.Printf("Received status of %s (%s)\n", monit.Server.HttpAddress, monit.Server.Localhostname)
	queue <- &monit
}

func main() {
	flag.Parse()

	go bootstrap()

	useLocal := os.Getenv("USE_LOCAL_FS") == "true"
	if useLocal {
		log.Println("Using local assets from 'public' directory")
	}

	http.Handle("/", http.FileServer(FS(useLocal)))
	http.HandleFunc("/status", status)
	http.HandleFunc("/action", action)
	http.HandleFunc("/collector", collector)

	log.Printf("Starting dogpack HTTP server at :%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
