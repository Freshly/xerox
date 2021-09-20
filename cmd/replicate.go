package cmd

import "github.com/spf13/cobra"

// rootCmd represents the base command when called without any subcommands
var replicateCmd = &cobra.Command{
	Use:   "replicate",
	Short: "xerox replicate replicates a REDIS instance --target-instance to --source-instance",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var (
	sourceInstance string
	sourceProject  string
	targetInstance string
	targetProject  string
)

func init() {
	replicateCmd.Flags().StringVar(&sourceInstance, "source-instance", "", "source REDIS instance name")
	replicateCmd.Flags().StringVar(&sourceProject, "source-project", "", "source REDIS instance project")
	replicateCmd.Flags().StringVar(&targetInstance, "target-instance", "", "target REDIS instance name")
	replicateCmd.Flags().StringVar(&targetProject, "target-project", "", "target REDIS instance project")
	rootCmd.AddCommand(replicateCmd)
}
