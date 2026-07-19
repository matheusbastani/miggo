package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/matheusbastani/miggo/internal/commands"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd.AddCommand(commands.Commands...)

	if err := rootCmd.Execute(); err != nil {
		color.Red("%s", err)
	}
}

var rootCmd = &cobra.Command{
	Use:           "miggo",
	Short:         "miggo: a simple SQL migration tool",
	SilenceErrors: true,
	SilenceUsage:  true,

	Run: func(cmd *cobra.Command, args []string) {
		shell(cmd)
	},
}

func shell(cmd *cobra.Command) {
	color.Yellow("  __  __ _               ")
	color.Yellow(" |  \\/  (_)__ _ __ _ ___ ")
	color.Yellow(" | |\\/| | / _` / _` / _ \\")
	color.Yellow(" |_|  |_|_\\__, \\__, \\___/")
	color.Yellow("          |___/|___/    ")
	color.Yellow("")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("miggo> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		args := strings.Fields(input)

		cmd.SetArgs(args)

		if err := cmd.Execute(); err != nil {
			color.Red("%s", err)
		}

		cmd.SetArgs(nil)
	}

	if err := scanner.Err(); err != nil {
		color.Red("%s", err)
	}
}
