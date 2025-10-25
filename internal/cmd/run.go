package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:   "run [name]",
	Short: "Run the main application logic",
	Long: `Run executes the main application logic with the provided configuration.
You can specify a name as an argument, or it will use the name from the config file.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := cfg.App.Name
		if len(args) > 0 {
			name = args[0]
		}

		message := cfg.App.Message
		if message == "" {
			message = "Hello"
		}

		result := map[string]interface{}{
			"name":    name,
			"message": fmt.Sprintf("%s, %s!", message, name),
			"timeout": cfg.App.Timeout,
			"debug":   cfg.App.Debug,
		}

		outputFormat := viper.GetString("output")
		verbose := viper.GetBool("verbose")

		if outputFormat == "json" {
			jsonData, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(jsonData))
		} else {
			fmt.Printf("%s\n", result["message"])
			if verbose {
				fmt.Printf("Timeout: %ds\n", result["timeout"])
				fmt.Printf("Debug: %v\n", result["debug"])
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("message", "m", "", "custom message to display")
	viper.BindPFlag("app.message", runCmd.Flags().Lookup("message"))
}
