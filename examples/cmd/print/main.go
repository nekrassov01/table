// Command print writes selected table examples.
package main

import (
	"fmt"
	"os"

	"github.com/nekrassov01/table/examples"
)

func main() {
	if err := examples.Run(os.Stdout, os.Args[1:]...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
