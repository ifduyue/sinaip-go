package main

import (
	"flag"
	"fmt"
)

func queryCmd() command {
	fs := flag.NewFlagSet("sinaip-go query", flag.ExitOnError)
	opts := &queryOpts{}

	return command{fs, "[options] <ip>...", func(args []string) error {
		fs.Parse(args)
		opts.ips = fs.Args()
		return query(opts)
	}}
}

type queryOpts struct {
	ips []string
}

func query(opts *queryOpts) error {
	for _, ip := range opts.ips {
		result, err := sinaip.Query(ip)
		if err != nil {
			fmt.Printf("%s:\t%v\n", ip, err)
		} else {
			fmt.Printf("%s:\t%s\n", ip, result.JSON())
		}
	}
	return nil
}
