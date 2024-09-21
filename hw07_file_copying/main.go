package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	parseFlags()
	if err := Copy(from, to, offset, limit); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func parseFlags() {
	flag.Parse()
	validateStringNotEmpty("file to read from", from)
	validateStringNotEmpty("file to write to", to)
	validatePositiveInt("limit", limit)
	validatePositiveInt("offset", offset)
}

func validateStringNotEmpty(name, value string) {
	if value == "" {
		fmt.Printf("%s is required\n", name)
		os.Exit(1)
	}
}

func validatePositiveInt(name string, value int64) {
	if value < 0 {
		fmt.Printf("%s must be positive\n", name)
		os.Exit(1)
	}
}
