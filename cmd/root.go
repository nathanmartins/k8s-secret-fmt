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
	var secret Secret
	err := yaml.Unmarshal(input, &secret)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %v", err)
	}

	var out strings.Builder
	encoder := yaml.NewEncoder(&out)
	encoder.SetIndent(0)

	err = encoder.Encode(secret)
	if err != nil {
		return nil, fmt.Errorf("error encoding YAML: %v", err)
	}

	result := strings.Builder{}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))

	inStringData := false
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "stringData:") {
			inStringData = true
			result.WriteString(line + "\n")
			continue
		}

		if inStringData && strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				value := strings.TrimSpace(parts[1])
				if !strings.HasPrefix(value, "'") {
					parts[1] = " '" + strings.Trim(value, "'\"") + "'"
				}
				line = strings.Join(parts, ":")
			}
		}

		result.WriteString(line + "\n")

		if inStringData && strings.TrimSpace(line) == "" {
			inStringData = false
		}
	}

	return []byte(result.String()), nil
}
