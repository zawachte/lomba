package cmd

import (
	"flag"
	"os"

	"github.com/spf13/cobra"
)

type Options struct {
}

var opts = &Options{}

var RootCmd = &cobra.Command{
	Use:          "lomba",
	SilenceUsage: true,
	Short:        "\n",
	Long:         "",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		// TODO: print error stack if log v>0
		// TODO: print cmd help if validation error
		os.Exit(1)
	}
}

func init() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
