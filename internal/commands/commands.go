package commands

import (
	"os"
	"strconv"

	"github.com/matheusbastani/miggo/internal/migrations"
	"github.com/matheusbastani/miggo/internal/settings"
	"github.com/spf13/cobra"
)

var Commands = []*cobra.Command{
	initCmd,
	createCmd,
	versionCmd,
	upCmd,
	downCmd,
	lockCmd,
	unlockCmd,
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
		return settings.CreateSettingsYAML()
	},
}

var createCmd = &cobra.Command{
	Use:   "create [name] [database]",
	Short: "Create a new migration",
	Long:  "Create creates a new migration with the given name",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, set, err := settings.GetDatabase(args[1])
		if err != nil {
			return err
		}

		return migrations.Create(set.Path, args[0])
	},
}

var versionCmd = &cobra.Command{
	Use:   "version [database]",
	Short: "Display the latest applied migration",
	Long:  "Display the latest applied migration folder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, _, err := getDatabase(args[0])
		if err != nil {
			return err
		}

		defer closeDatabase(db)

		return migrations.Version(db)
	},
}

var upCmd = &cobra.Command{
	Use:   "up [database]",
	Short: "Apply all pending migrations",
	Long:  "Up applies all pending migrations to the database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, set, err := getDatabase(args[0])
		if err != nil {
			return err
		}

		defer closeDatabase(db)

		return migrations.Up(db, set.Path)
	},
}

var downCmd = &cobra.Command{
	Use:   "down [database]",
	Short: "Rollback the most recently applied migration",
	Long:  "Down rolls back the most recently applied migration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, set, err := getDatabase(args[0])
		if err != nil {
			return err
		}

		defer closeDatabase(db)

		return migrations.Down(db, set.Path)
	},
}

var lockCmd = &cobra.Command{
	Use:   "lock [index] [database]",
	Short: "Lock a migration",
	Long:  "Lock a migration, preventing it from being rolled back",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, _, err := getDatabase(args[1])
		if err != nil {
			return err
		}

		defer closeDatabase(db)

		return migrations.SetRollbackBoundary(db, args[0])
	},
}

var unlockCmd = &cobra.Command{
	Use:   "unlock [index] [database]",
	Short: "Unlock a migration",
	Long:  "Unlock a migration, allowing it to be rolled back",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, _, err := getDatabase(args[1])
		if err != nil {
			return err
		}

		defer closeDatabase(db)

		return migrations.RemoveRollbackBoundary(db, args[0])
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset [database]",
	Short: "Rollback all applied migrations",
	Long:  "Reset rolls back all applied migrations in reverse order",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, set, err := getDatabase(args[0])
		if err != nil {
			return err
		}

		defer closeDatabase(db)

		force, err := getForce(cmd)
		if err != nil {
			return err
		}

		secure, err := getSecure(set.Environment)
		if err != nil {
			return err
		}

		return migrations.Reset(
			db,
			set.Path,
			secure,
			force,
		)
	},
}

var resetDropCmd = &cobra.Command{
	Use:   "reset-drop [database]",
	Short: "Rollback all migrations and drop the migrations table",
	Long:  "Reset and drop rolls back all migrations and drops the miggo table",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, set, err := getDatabase(args[0])
		if err != nil {
			return err
		}

		defer closeDatabase(db)

		force, err := getForce(cmd)
		if err != nil {
			return err
		}

		secure, err := getSecure(set.Environment)
		if err != nil {
			return err
		}

		return migrations.ResetAndDrop(
			db,
			set.Path,
			secure,
			force,
		)
	},
}

var insertCmd = &cobra.Command{
	Use:   "insert [name] [index] [database]",
	Short: "Create a new migration at a specific index",
	Long:  "Insert creates a new migration at a specific index, renumbering existing migrations as needed",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, set, err := getDatabase(args[2])
		if err != nil {
			return err
		}

		defer closeDatabase(db)

		index, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}

		force, err := getForce(cmd)
		if err != nil {
			return err
		}

		secure, err := getSecure(set.Environment)
		if err != nil {
			return err
		}

		return migrations.Insert(
			db,
			set.Path,
			args[0],
			index,
			secure,
			force,
		)
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
