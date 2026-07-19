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
		_, settings, err := settings.GetDatabase(args[1])
		if err != nil {
			return err
		}

		return migrations.Create(settings.Path, args[0])
	},
}

var versionCmd = &cobra.Command{
	Use:   "version [database]",
	Short: "Display the latest applied migration",
	Long:  "Display the latest applied migration folder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, _, err := settings.GetDatabase(args[0])
		if err != nil {
			return err
		}

		defer func() {
			_ = db.Close()
		}()

		return migrations.Version(db)
	},
}

var upCmd = &cobra.Command{
	Use:   "up [database]",
	Short: "Apply all pending migrations",
	Long:  "Up applies all pending migrations to the database. It creates the miggo table if it doesn't exist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, settings, err := settings.GetDatabase(args[0])
		if err != nil {
			return err
		}

		defer func() {
			_ = db.Close()
		}()

		return migrations.Up(db, settings.Path)
	},
}

var downCmd = &cobra.Command{
	Use:   "down [database]",
	Short: "Rollback the most recently applied migration",
	Long:  "Down rolls back the most recently applied migration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, settings, err := settings.GetDatabase(args[0])
		if err != nil {
			return err
		}

		defer func() {
			_ = db.Close()
		}()

		return migrations.Down(db, settings.Path)
	},
}

var lockCmd = &cobra.Command{
	Use:   "lock [index] [database]",
	Short: "Lock a migration",
	Long:  "Lock a migration, preventing it from being rolled back",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, _, err := settings.GetDatabase(args[1])
		if err != nil {
			return err
		}

		migration := args[0]

		return migrations.SetRollbackBoundary(db, migration)
	},
}

var unlockCmd = &cobra.Command{
	Use:   "unlock [index] [database]",
	Short: "Unlock a migration",
	Long:  "Unlock a migration, allowing it to be rolled back",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, _, err := settings.GetDatabase(args[1])
		if err != nil {
			return err
		}

		migration := args[0]

		return migrations.RemoveRollbackBoundary(db, migration)
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset [database]",
	Short: "Rollback all applied migrations",
	Long:  "Reset rolls back all applied migrations in reverse order",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, settings, err := settings.GetDatabase(args[0])
		if err != nil {
			return err
		}

		defer func() {
			_ = db.Close()
		}()

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		return migrations.Reset(
			db,
			settings.Path,
			settings.Secure,
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
		db, settings, err := settings.GetDatabase(args[0])
		if err != nil {
			return err
		}

		defer func() {
			_ = db.Close()
		}()

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		return migrations.ResetAndDrop(
			db,
			settings.Path,
			settings.Secure,
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
		db, settings, err := settings.GetDatabase(args[2])
		if err != nil {
			return err
		}

		migration := args[0]

		index, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		return migrations.Insert(db, settings.Path, migration, index, settings.Secure, force)
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
