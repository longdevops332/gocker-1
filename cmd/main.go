package cmd

import (
	"os"

	"github.com/saromanov/gocker/pkg/cmd"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)


func main() {
	cmd.Build(os.Args)
}
