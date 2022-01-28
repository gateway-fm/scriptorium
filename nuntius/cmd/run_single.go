package cmd

import (
	"github.com/gateway-fm/scriptorium/nuntius/internal"
	"github.com/spf13/cobra"
)

// Cmd is
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run-single",
		Short: "Run single request",
		//PersistentPreRun: func(cmd *cobra.Command, args []string) {
		//	payload := &config.Config{}
		//	err := payload.ParsePayload()
		//	if err != nil{
		//		logger.Log().Error(fmt.Errorf("parsing payload failed: %s", err).Error())
		//	}
		//},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := internal.TempRequest(); err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}
