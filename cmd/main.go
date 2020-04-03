package main

import (
	"os"

	"github.com/saromanov/gocker/pkg/cmd"
)


func main() {
	cmd.Build(os.Args)
}
