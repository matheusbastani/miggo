package commands

import (
	"database/sql"

	"github.com/matheusbastani/miggo/internal/settings"
	"github.com/spf13/cobra"
)

func getDatabase(name string) (*sql.DB, settings.Database, error) {
	db, set, err := settings.GetDatabase(name)
	if err != nil {
		return nil, settings.Database{}, err
	}

	return db, set, nil
}

func closeDatabase(db *sql.DB) {
	_ = db.Close()
}

func getForce(cmd *cobra.Command) (bool, error) {
	return cmd.Flags().GetBool("force")
}

func getSecure(environment settings.Environment) (bool, error) {
	return settings.IsSecure(environment)
}
