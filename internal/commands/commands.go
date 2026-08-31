package commands

import (
	"os"
	"strconv"

	"github.com/matheusbastani/miggo/internal/migrations"
	"github.com/matheusbastani/miggo/internal/settings"
	"github.com/spf13/cobra"
)

// NewCommands returns a fresh set of command instances.
//
// A new instance is required every call because cobra commands carry
// flag state between executions; reusing the same *cobra.Command across
// multiple Execute() calls (as the shell does) leaks flag values between
// invocations.
func NewCommands() []*cobra.Command {
	return []*cobra.Command{
		newInitCmd(),
		newCreateCmd(),
		newVersionCmd(),
		newUpCmd(),
		newDownCmd(),
		newLockCmd(),
		newUnlockCmd(),
		newResetCmd(),
		newResetDropCmd(),
		newInsertCmd(),
		NewExitCmd(),
	}
}

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create a miggo.yaml file",
		Long:  "Create a miggo.yaml file in the current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			return settings.CreateSettingsYAML()
		},
	}
}

func newCreateCmd() *cobra.Command {
	var dbName string

	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new migration",
		Long:  "Create creates a new migration with the given name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, set, err := settings.GetDatabase(dbName)
			if err != nil {
				return err
			}

			return migrations.Create(set.Path, args[0])
		},
	}

	cmd.Flags().StringVarP(&dbName, "db", "d", "", "database to use")
	_ = cmd.MarkFlagRequired("db")

	return cmd
}

func newVersionCmd() *cobra.Command {
	var dbName string

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display the latest applied migration",
		Long:  "Display the latest applied migration folder",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			db, _, err := getDatabase(dbName)
			if err != nil {
				return err
			}

			defer closeDatabase(db)

			return migrations.Version(db)
		},
	}

	cmd.Flags().StringVarP(&dbName, "db", "d", "", "database to use")
	_ = cmd.MarkFlagRequired("db")

	return cmd
}

func newUpCmd() *cobra.Command {
	var dbName string

	cmd := &cobra.Command{
		Use:   "up",
		Short: "Apply all pending migrations",
		Long:  "Up applies all pending migrations to the database",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			db, set, err := getDatabase(dbName)
			if err != nil {
				return err
			}

			defer closeDatabase(db)

			return migrations.Up(db, set.Driver, set.Path)
		},
	}

	cmd.Flags().StringVarP(&dbName, "db", "d", "", "database to use")
	_ = cmd.MarkFlagRequired("db")

	return cmd
}

func newDownCmd() *cobra.Command {
	var dbName string

	cmd := &cobra.Command{
		Use:   "down",
		Short: "Rollback the most recently applied migration",
		Long:  "Down rolls back the most recently applied migration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			db, set, err := getDatabase(dbName)
			if err != nil {
				return err
			}

			defer closeDatabase(db)

			return migrations.Down(db, set.Driver, set.Path)
		},
	}

	cmd.Flags().StringVarP(&dbName, "db", "d", "", "database to use")
	_ = cmd.MarkFlagRequired("db")

	return cmd
}

func newLockCmd() *cobra.Command {
	var dbName string

	cmd := &cobra.Command{
		Use:   "lock [index]",
		Short: "Lock a migration",
		Long:  "Lock a migration, preventing it from being rolled back",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, set, err := getDatabase(dbName)
			if err != nil {
				return err
			}

			defer closeDatabase(db)

			return migrations.SetRollbackBoundary(db, set.Driver, args[0])
		},
	}

	cmd.Flags().StringVarP(&dbName, "db", "d", "", "database to use")
	_ = cmd.MarkFlagRequired("db")

	return cmd
}

func newUnlockCmd() *cobra.Command {
	var dbName string

	cmd := &cobra.Command{
		Use:   "unlock [index]",
		Short: "Unlock a migration",
		Long:  "Unlock a migration, allowing it to be rolled back",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, set, err := getDatabase(dbName)
			if err != nil {
				return err
			}

			defer closeDatabase(db)

			return migrations.RemoveRollbackBoundary(db, set.Driver, args[0])
		},
	}

	cmd.Flags().StringVarP(&dbName, "db", "d", "", "database to use")
	_ = cmd.MarkFlagRequired("db")

	return cmd
}

func newResetCmd() *cobra.Command {
	var dbName string

	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Rollback all applied migrations",
		Long:  "Reset rolls back all applied migrations in reverse order",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			db, set, err := getDatabase(dbName)
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
				set.Driver,
				set.Path,
				secure,
				force,
			)
		},
	}

	cmd.Flags().StringVarP(&dbName, "db", "d", "", "database to use")
	_ = cmd.MarkFlagRequired("db")

	return cmd
}

func newResetDropCmd() *cobra.Command {
	var dbName string

	cmd := &cobra.Command{
		Use:   "reset-drop",
		Short: "Rollback all migrations and drop the migrations table",
		Long:  "Reset and drop rolls back all migrations and drops the miggo table",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			db, set, err := getDatabase(dbName)
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
				set.Driver,
				set.Path,
				secure,
				force,
			)
		},
	}

	cmd.Flags().StringVarP(&dbName, "db", "d", "", "database to use")
	_ = cmd.MarkFlagRequired("db")

	return cmd
}

func newInsertCmd() *cobra.Command {
	var dbName string

	cmd := &cobra.Command{
		Use:   "insert [name] [index]",
		Short: "Create a new migration at a specific index",
		Long:  "Insert creates a new migration at a specific index, renumbering existing migrations as needed",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, set, err := getDatabase(dbName)
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

	cmd.Flags().StringVarP(&dbName, "db", "d", "", "database to use")
	_ = cmd.MarkFlagRequired("db")

	return cmd
}

// NewExitCmd returns the "exit" command, which exits miggo.
func NewExitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "exit",
		Short: "Exit miggo",
		Long:  "Exit miggo",
		Run: func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		},
	}
}
