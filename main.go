package main

import (
	"flag"
	"fmt"
	sinaiplib "github.com/ifduyue/sinaip-go/lib"
	"log"
	"os"
	"runtime"
)

var (
	sinaip *sinaiplib.SINAIP
)

func main() {
	commands := map[string]command{"httpd": httpdCmd(), "query": queryCmd()}

	flag.Usage = func() {
		fmt.Println("Usage: sinaip-go [globals] <command> [options]")
		for name, cmd := range commands {
			fmt.Printf("\n%s command: %s\n", name, cmd.usage)
			cmd.fs.PrintDefaults()
		}
		fmt.Printf("\nglobal flags:\n")
		fmt.Printf("\t-cpus=%d Number of CPUs to use\n", runtime.NumCPU())
		fmt.Printf("\t-ipdat=\"\" Path to ip.dat, will try to get it from env variable \"SINAIPDAT\" if left empty.\n")
		fmt.Println(examples)
	}

	cpus := flag.Int("cpus", runtime.NumCPU(), "Number of CPUs to use")
	ipdatpath := flag.String("ipdat", os.Getenv("SINAIPDAT"), "Path to ip.dat")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	runtime.GOMAXPROCS(*cpus)
	if *ipdatpath != "" {
		os.Setenv("SINAIPDAT", *ipdatpath)
	}

	if os.Getenv("SINAIPDAT") == "" {
		log.Fatal("Path to ip.dat can't be empty")
	}

	if cmd, ok := commands[args[0]]; !ok {
		log.Fatalf("Unknown command: %s", args[0])
	} else if err := cmd.fn(args[1:]); err != nil {
		log.Fatal(err)
	}
}

const examples = `
examples:
  sinaip-go query 1.2.3.4
  sinaip-go httpd 127.0.0.1:8080
`

type command struct {
	fs    *flag.FlagSet
	usage string
	fn    func(args []string) error
}
