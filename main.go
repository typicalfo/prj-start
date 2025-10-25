package main

import (
	"github.com/typicalfo/prj-start/cmd"
	"github.com/typicalfo/prj-start/logger"
)

func main() {
	logger.InitLogger()
	cmd.Execute()
}
