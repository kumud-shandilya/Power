package main

import (
	"github.com/getporter/power/pkg/power"
	"github.com/spf13/cobra"
)

func buildSchemaCommand(m *power.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Print the json schema for the mixin",
		Run: func(cmd *cobra.Command, args []string) {
			m.PrintSchema()
		},
	}
	return cmd
}
