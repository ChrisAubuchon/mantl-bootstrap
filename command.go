package main

import (
	"github.com/asteris-llc/mantl-bootstrap/bootstrap"
	"github.com/asteris-llc/mantl-bootstrap/listen"

	"github.com/spf13/cobra"
)

func initCommand(name, version string) *cobra.Command {
	root := &cobra.Command{
		Use:   "mantl-bootstrap",
		Short: "Bootstrap a list of nodes with mantl",
		Long:  "Bootstrap a list of nodes with mantl",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	bootstrap.Init(root)
	listen.Init(root)

	return root
}
