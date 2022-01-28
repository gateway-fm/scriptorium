package main

import (
	"fmt"
	"os"

	"github.com/gateway-fm/scriptorium/logger"
	"github.com/gateway-fm/scriptorium/nuntius/cmd"
)

func main() {

	singlerequst := cmd.Cmd()
	singlerequst.AddCommand(cmd.Cmd())
	if err := singlerequst.Execute(); err != nil {
		logger.Log().Error(fmt.Errorf("single request failed: %w", err).Error())
		os.Exit(1)
	}
}
