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
		fmt.Println("Validation error:", err)
		os.Exit(1)
	}

	if err := Copy(from, to, offset, limit); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func validateArgs() error {
	validations := []struct {
		name  string
		value interface{}
		fn    func(string, interface{}) error
	}{
		{"file to read from", from, validateStringNotEmpty},
		{"file to write to", to, validateStringNotEmpty},
		{"limit", limit, validatePositiveInt},
		{"offset", offset, validatePositiveInt},
	}

	for _, v := range validations {
		if err := v.fn(v.name, v.value); err != nil {
			return err
		}
	}

	return nil
}

func validateStringNotEmpty(name string, value interface{}) error {
	if value == "" {
		return fmt.Errorf("%s is required", name)
	}
	return nil
}

func validatePositiveInt(name string, value interface{}) error {
	v, ok := value.(int64)
	if !ok {
		return fmt.Errorf("%s must be an integer", name)
	}
	if v < 0 {
		return fmt.Errorf("%s must be positive", name)
	}
	return nil
}
