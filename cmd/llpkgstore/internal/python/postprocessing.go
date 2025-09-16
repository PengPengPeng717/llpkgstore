package internal

import (
	"github.com/goplus/llpkgstore/internal/actions"
	"github.com/spf13/cobra"
)

var postProcessingCmd = &cobra.Command{
	Use:   "postprocessing",
	Short: "Process merged PR for Python packages",
	Long:  `Process merged PR for Python packages with unified version management`,
	RunE:  runPythonPostProcessingCmd,
}

func runPythonPostProcessingCmd(_ *cobra.Command, _ []string) error {
	client, err := actions.NewDefaultClient()
	if err != nil {
		return err
	}
	return client.Postprocessing()
}

func init() {
	rootCmd.AddCommand(postProcessingCmd)
}
