package main

import (
	"github.com/getporter/power/pkg/power"
	"github.com/spf13/cobra"
)

func buildUpgradeCommand(m *power.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Execute the invoke functionality of this mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Execute(cmd.Context())
		},
	}
	return cmd
}
