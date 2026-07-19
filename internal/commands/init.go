package commands

func init() {
	resetCmd.Flags().Bool(
		"force",
		false,
		"ignore rollback boundaries",
	)

	resetDropCmd.Flags().Bool(
		"force",
		false,
		"ignore rollback boundaries",
	)
	insertCmd.Flags().Bool(
		"force",
		false,
		"ignore rollback boundaries",
	)
}
