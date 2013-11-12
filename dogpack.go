package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	port       = flag.String("p", "2813", "Port to serve dogpack")
	public_dir = flag.String("d", "./public", "Public directory location")
)

type DogPack map[string]*Monit

var (
	dogpack = &DogPack{}
	queue   = make(chan *Monit)
)

func bootstrap() {
	for {
		monit := <-queue
		(*dogpack)[monit.Server.Localhostname] = monit
	}
}

func Collector(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	monit := MonitFromXml(req.Body)
	monit.IpAddress = strings.Split(req.RemoteAddr, ":")[0]
	queue <- &monit
}

func Status(w http.ResponseWriter, req *http.Request) {
	json, err := json.Marshal(dogpack)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(json)
	}
}

func Action(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	if req.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	server := req.PostFormValue("server")
	service := req.PostFormValue("service")
	action := req.PostFormValue("action")

	monit, found := (*dogpack)[server]
	if found {
		log.Printf("Server[%s]: Executing action '%s' on service '%s'", server, action, service)
		go monit.Server.Execute(service, action)
	} else {
		log.Printf("Server[%s]: Unable to execute action '%s' on service '%s'. Server not found.", server, action, service)
	}
}

func main() {
	flag.Parse()

	go bootstrap()

	http.HandleFunc("/collector", Collector)
	http.HandleFunc("/status", Status)
	http.HandleFunc("/action", Action)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(*public_dir))))

	address := fmt.Sprintf(":%s", *port)
	log.Println("Starting dogpack HTTP server at ", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
