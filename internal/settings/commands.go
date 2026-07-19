package settings

import (
	"os"
	"strconv"

	"github.com/matheusbastani/miggo/internal/migrations"
	"github.com/spf13/cobra"
)

var Commands = []*cobra.Command{
	initCmd,
	createCmd,
	versionCmd,
	upCmd,
	downCmd,
	resetCmd,
	resetDropCmd,
	insertCmd,
	exitCmd,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a miggo.yaml file",
	Long:  "Create a miggo.yaml file in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		return CreateSettingsYAML()
	},
}

var createCmd = &cobra.Command{
	Use:   "create [name] [database]",
	Short: "Create a new migration",
	Long:  "Create creates a new migration with the given name",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, path, err := GetDatabase(args[1])
		if err != nil {
			return err
		}

		return migrations.Create(path, args[0])
	},
}

var versionCmd = &cobra.Command{
	Use:   "version [database]",
	Short: "Display the latest applied migration",
	Long:  "Display the latest applied migration folder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, _, err := GetDatabase(args[0])
		if err != nil {
			return err
		}

		defer db.Close()

		return migrations.Version(db)
	},
}

var upCmd = &cobra.Command{
	Use:   "up [database]",
	Short: "Apply all pending migrations",
	Long:  "Up applies all pending migrations to the database. It creates the miggo table if it doesn't exist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, path, err := GetDatabase(args[0])
		if err != nil {
			return err
		}

		defer db.Close()

		return migrations.Up(db, path)
	},
}

var downCmd = &cobra.Command{
	Use:   "down [database]",
	Short: "Rollback the most recently applied migration",
	Long:  "Down rolls back the most recently applied migration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, path, err := GetDatabase(args[0])
		if err != nil {
			return err
		}

		defer db.Close()

		return migrations.Down(db, path)
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset [database]",
	Short: "Rollback all applied migrations",
	Long:  "Reset rolls back all applied migrations in reverse order",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, path, err := GetDatabase(args[0])
		if err != nil {
			return err
		}

		defer db.Close()

		return migrations.Reset(db, path)
	},
}

var resetDropCmd = &cobra.Command{
	Use:   "reset-drop [database]",
	Short: "Rollback all migrations and drop the migrations table",
	Long:  "Reset and drop rolls back all migrations and drops the miggo table",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, path, err := GetDatabase(args[0])
		if err != nil {
			return err
		}

		defer db.Close()

		return migrations.ResetAndDrop(db, path)
	},
}

var insertCmd = &cobra.Command{
	Use:   "insert [name] [index] [database]",
	Short: "Create a new migration at a specific index",
	Long:  "Insert creates a new migration at a specific index, renumbering existing migrations as needed",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, path, err := GetDatabase(args[2])
		if err != nil {
			return err
		}

		index, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}

		return migrations.Insert(path, args[0], index)
	},
}

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Exit miggo",
	Long:  "Exit miggo",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}
