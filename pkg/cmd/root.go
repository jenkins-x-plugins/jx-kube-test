package cmd

import (
	"github.com/jenkins-x-plugins/jx-kube-test/pkg/cmd/run"
	"github.com/jenkins-x-plugins/jx-kube-test/pkg/cmd/version"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/spf13/cobra"
)

// Main creates the new command
func Main() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kube-test",
		Short: "commands for working with GitOps based git repositories",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				log.Logger().Errorf(err.Error())
			}
		},
	}
	cmd.AddCommand(cobras.SplitCommand(run.NewCmdRun()))
	cmd.AddCommand(cobras.SplitCommand(version.NewCmdVersion()))
	return cmd
}
