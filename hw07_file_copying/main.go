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
	flag.Parse()

	if err := validateArgs(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if err := Copy(from, to, offset, limit); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func validateArgs() error {
	if err := validateStringNotEmpty("file to read from", from); err != nil {
		return err
	}
	if err := validateStringNotEmpty("file to write to", to); err != nil {
		return err
	}
	if err := validatePositiveInt("limit", limit); err != nil {
		return err
	}
	if err := validatePositiveInt("offset", offset); err != nil {
		return err
	}
	return nil
}

func validateStringNotEmpty(name, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", name)
	}
	return nil
}

func validatePositiveInt(name string, value int64) error {
	if value < 0 {
		return fmt.Errorf("%s must be positive", name)
	}
	return nil
}
