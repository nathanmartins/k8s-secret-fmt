package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessYAML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Simple Secret",
			input: `apiVersion: v1
kind: Secret
metadata:
  name: my-secret
type: Opaque
stringData:
  username: admin
  api-key: 1234-5678-abcd
`,
			expected: `apiVersion: v1
kind: Secret
metadata:
  name: my-secret
type: Opaque
stringData:
  api-key: '1234-5678-abcd'
  username: 'admin'
`,
		},
		{
			name: "Secret with keys in non-alphabetical order",
			input: `apiVersion: v1
kind: Secret
metadata:
  name: my-secret
type: Opaque
stringData:
  zebra: striped
  apple: red
  banana: yellow
`,
			expected: `apiVersion: v1
kind: Secret
metadata:
  name: my-secret
type: Opaque
stringData:
  apple: 'red'
  banana: 'yellow'
  zebra: 'striped'
`,
		},
		{
			name: "Secret with multi-line value",
			input: `apiVersion: v1
kind: Secret
metadata:
 annotations:
  kustomize.config.k8s.io/needs-hash: "true"
  name: example
type: Opaque
stringData:
 config.yaml: |
   app:
    env: production
    host: 0.0.0.0
		`,
			expected: `apiVersion: v1
kind: Secret
metadata:
 annotations:
  kustomize.config.k8s.io/needs-hash: "true"
  name: example
type: Opaque
stringData:
 config.yaml: |
   app:
    env: production
    host: 0.0.0.0
		`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := processYAML([]byte(tt.input))
			if err != nil {
				t.Fatalf("processYAML() error = %v", err)
			}

			// Normalize line endings
			outputStr := string(output)
			expectedStr := tt.expected

			assert.Equal(t, strings.TrimSpace(expectedStr), strings.TrimSpace(outputStr))
		})
	}
}
