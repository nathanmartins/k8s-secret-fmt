package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "k8s-secret-fmt",
	Run: func(cmd *cobra.Command, args []string) {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}

		output, err := processYAML(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing YAML: %v\n", err)
			os.Exit(1)
		}

		fmt.Print(string(output))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.k8s-secret-fmt.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Secret struct {
	ApiVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Metadata   map[string]interface{} `yaml:"metadata"`
	Type       string                 `yaml:"type"`
	StringData map[string]string      `yaml:"stringData,omitempty"`
}

func processYAML(input []byte) ([]byte, error) {
	// Parse the original YAML to get the structure
	var secret Secret
	err := yaml.Unmarshal(input, &secret)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %v", err)
	}

	// Create a new YAML with proper indentation
	var out strings.Builder
	encoder := yaml.NewEncoder(&out)
	encoder.SetIndent(2)

	// Create a modified secret with the same data but without stringData
	// We'll handle stringData separately to ensure proper formatting
	modifiedSecret := Secret{
		ApiVersion: secret.ApiVersion,
		Kind:       secret.Kind,
		Metadata:   secret.Metadata,
		Type:       secret.Type,
		// Omit StringData as we'll handle it separately
	}

	err = encoder.Encode(modifiedSecret)
	if err != nil {
		return nil, fmt.Errorf("error encoding YAML: %v", err)
	}

	// Process the YAML line by line
	result := strings.Builder{}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))

	for scanner.Scan() {
		line := scanner.Text()
		result.WriteString(line + "\n")
	}

	// Add stringData section if it exists
	if len(secret.StringData) > 0 {
		result.WriteString("stringData:\n")

		// Process all stringData fields and add them with proper formatting
		for key, value := range secret.StringData {
			// Handle multi-line values
			if strings.Contains(value, "\n") {
				// Use YAML literal block style for multi-line values
				result.WriteString("  " + key + ": |-\n")
				scanner = bufio.NewScanner(strings.NewReader(value))
				for scanner.Scan() {
					result.WriteString("    " + scanner.Text() + "\n")
				}
			} else {
				// Use single quotes for single-line values
				result.WriteString("  " + key + ": '" + value + "'\n")
			}
		}
	}

	return []byte(result.String()), nil
}
