package main

import (
	"errors"
	"flag"
	sinaiplib "github.com/ifduyue/sinaip-go/lib"
	"log"
	"net/http"
	"os"
)

func httpdCmd() command {
	fs := flag.NewFlagSet("sinaip-go httpd", flag.ExitOnError)
	opts := &httpdOpts{}

	return command{fs, "[options] <addr>", func(args []string) error {
		fs.Parse(args)
		opts.addrs = fs.Args()
		if len(opts.addrs) != 1 {
			return errors.New("You must specify a http addr.")
		}
		return httpd(opts)
	}}
}

type httpdOpts struct {
	addrs []string
}

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

func httpd(opts *httpdOpts) error {
	var err error
	sinaip, err = sinaiplib.NewSINAIP(os.Getenv("SINAIPDAT"))
	if err != nil {
		return err
	}

	httpaddr := opts.addrs[0]
	log.Printf("%10v", sinaip.Size)
	log.Printf("%10v", sinaip.IndexOffset)
	log.Printf("%10v", sinaip.Count)
	log.Printf("Listening to %v", httpaddr)
	http.HandleFunc("/", rootHandler)
	return http.ListenAndServe(httpaddr, nil)
}
