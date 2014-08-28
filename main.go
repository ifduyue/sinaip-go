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
	sinaip    *sinaiplib.SINAIP
	cpus      *int
	ipdatpath *string
	preload   *bool
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
		fmt.Printf("\t-preload If true, preload ip.dat into memory, otherwise mmap is used.\n")
		fmt.Println(examples)
	}

	cpus = flag.Int("cpus", runtime.NumCPU(), "Number of CPUs to use")
	ipdatpath = flag.String("ipdat", os.Getenv("SINAIPDAT"), "Path to ip.dat")
	preload = flag.Bool("preload", false, "Preload ip.dat to memory, otherwise mmap is used.")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	runtime.GOMAXPROCS(*cpus)
	if *ipdatpath == "" {
		log.Fatal("Path to ip.dat can't be empty")
	}

	if sinaiptmp, err := sinaiplib.NewSINAIP(*ipdatpath, *preload); err != nil {
		log.Fatal(err)
	} else {
		sinaip = sinaiptmp
	}

	if cmd, ok := commands[args[0]]; !ok {
		log.Fatalf("Unknown command: %s", args[0])
	} else if err := cmd.fn(args[1:]); err != nil {
		log.Fatal(err)
	}
}

const examples = `
examples:
\tsinaip-go query 1.2.3.4
\tsinaip-go httpd 127.0.0.1:8080
`

type command struct {
	fs    *flag.FlagSet
	usage string
	fn    func(args []string) error
}
