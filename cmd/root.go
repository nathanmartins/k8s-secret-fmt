package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"
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
	// Instead of parsing and re-encoding the YAML, we'll preserve the original format
	// but just make the necessary modifications to handle the stringData section correctly

	// First, let's parse the YAML to get the structure
	// Preprocess the input to replace tabs with spaces
	inputStr := string(input)
	inputStr = strings.ReplaceAll(inputStr, "\t", "  ")

	var secret Secret
	err := yaml.Unmarshal([]byte(inputStr), &secret)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %v", err)
	}

	// Now, let's reconstruct the YAML preserving the original format
	// We'll use the original input as a template and just modify the stringData section

	// Split the input into lines
	lines := strings.Split(inputStr, "\n")

	// Find the stringData section
	stringDataIndex := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "stringData:" {
			stringDataIndex = i
			break
		}
	}

	// If there's no stringData section or no stringData in the parsed secret, return the original input
	if stringDataIndex == -1 || len(secret.StringData) == 0 {
		return []byte(inputStr), nil
	}

	// Keep everything up to and including the stringData line
	result := strings.Join(lines[:stringDataIndex+1], "\n")

	// Add each stringData entry with the correct format
	// Get all keys and sort them alphabetically
	keys := make([]string, 0, len(secret.StringData))
	for key := range secret.StringData {
		keys = append(keys, key)
	}
	// Sort the keys alphabetically
	sort.Strings(keys)

	// Get the indentation from the original input
	indentation := ""
	for i := stringDataIndex + 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Find the first non-whitespace character
		for j, c := range line {
			if c != ' ' && c != '\t' {
				indentation = line[:j]
				break
			}
		}

		if indentation != "" {
			break
		}
	}

	// If we couldn't determine the indentation, use a single space
	if indentation == "" {
		indentation = " "
	}

	// Iterate through the sorted keys
	for _, key := range keys {
		value := secret.StringData[key]

		// Add the key
		result += "\n" + indentation + key + ": "

		// Handle multi-line values
		if strings.Contains(value, "\n") {
			// Use YAML literal block style for multi-line values
			result += "|"

			// Add each line of the value with the correct indentation
			valueLines := strings.Split(value, "\n")
			for _, valueLine := range valueLines {
				result += "\n" + indentation + " " + indentation + valueLine
			}
		} else {
			// Use single quotes for single-line values
			result += "'" + value + "'"
		}
	}

	return []byte(result), nil
}
