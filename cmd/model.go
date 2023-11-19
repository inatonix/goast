package cmd

import (
	"github.com/spf13/cobra"
)

type Command interface {
	Initialize()
}

type CommandModel struct {
	Root *cobra.Command
}

type cobraParameters struct {
	Use   string
	Short string
	Long  string
}

func (c *CommandModel) Initialize(params cobraParameters) error {
	c.Root = &cobra.Command{
		Use:   params.Use,
		Short: params.Short,
		Long:  params.Long,
	}
	return nil
}
