package main

import (
	"flag"
	"fmt"
	sinaiplib "github.com/ifduyue/sinaip.go"
	"log"
	"net/http"
)

var (
	port      int
	ipdatPath string
	sinaip    *sinaiplib.SINAIP
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	ipstr := r.URL.Path[1:]
	ip, err := sinaip.Query(ipstr)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.Write(ip.Json())
}

func init() {
	flag.IntVar(&port, "port", 8080, "HTTP Server Port")
	flag.StringVar(&ipdatPath, "ipdat", "./ip.dat", "Path to ip.dat")
	flag.Parse()
}

func main() {
	var err error
	sinaip, err = sinaiplib.NewSINAIP(ipdatPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	httpAddr := fmt.Sprintf(":%v", port)
	log.Printf("%10v", sinaip.Size)
	log.Printf("%10v", sinaip.IndexOffset)
	log.Printf("%10v", sinaip.Count)
	log.Printf("Listening to %v", httpAddr)
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(httpAddr, nil))
}
