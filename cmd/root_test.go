package cmd

import (
	"strings"
	"testing"
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
  password: complex"password'with"quotes
  api-key: 1234-5678-abcd
`,
			expected: `apiVersion: v1
kind: Secret
metadata:
  name: my-secret
type: Opaque
stringData:
  username: 'admin'
  password: 'complex"password''with"quotes'
  api-key: '1234-5678-abcd'
`,
		},
		{
			name: "Secret with multi-line value",
			input: `apiVersion: v1
kind: Secret
metadata:
    annotations:
        kustomize.config.k8s.io/needs-hash: "true"
    name: income-bureau-etl-secrets
type: Opaque
stringData:
    config.yaml: |
        app:
          env: production
          host: 0.0.0.0
          port: 50051
          rules_path: 'aggregation-rules.yaml'
`,
			expected: `apiVersion: v1
kind: Secret
metadata:
  annotations:
    kustomize.config.k8s.io/needs-hash: "true"
  name: income-bureau-etl-secrets
type: Opaque
stringData:
  config.yaml: |
    app:
      env: production
      host: 0.0.0.0
      port: 50051
      rules_path: 'aggregation-rules.yaml'
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
			outputStr := strings.ReplaceAll(string(output), "\r\n", "\n")
			expectedStr := strings.ReplaceAll(tt.expected, "\r\n", "\n")

			if outputStr != expectedStr {
				t.Errorf("processYAML() output = %v, want %v", outputStr, expectedStr)
			}
		})
	}
}
