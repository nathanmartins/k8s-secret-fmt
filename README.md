# k8s-secret-fmt

A command-line tool that formats Kubernetes secrets YAML files by ensuring string values in `stringData` fields are properly quoted. This is particularly useful when working with encrypted secrets using tools like [SOPS](https://github.com/mozilla/sops).

## Purpose

When working with Kubernetes secrets, especially when using encryption tools like SOPS, you might encounter issues with string formatting in the YAML files. This tool ensures that all values under the `stringData` field are properly single-quoted, preventing common YAML parsing issues.

## Installation

You can install k8s-secret-fmt using one of the following methods:

### Download Binary

Download the latest release from the [releases page](https://github.com/nathanmartins/k8s-secret-fmt/releases/latest) for your operating system:

```bash
# For macOS (x86_64)
curl -L https://github.com/nathanmartins/k8s-secret-fmt/releases/download/v1.0.0/k8s-secret-fmt_Darwin_x86_64.tar.gz | tar xz
sudo mv k8s-secret-fmt /usr/local/bin/

# For macOS (Apple Silicon/ARM64)
curl -L https://github.com/nathanmartins/k8s-secret-fmt/releases/download/v1.0.0/k8s-secret-fmt_Darwin_arm64.tar.gz | tar xz
sudo mv k8s-secret-fmt /usr/local/bin/

# For Linux (x86_64)
curl -L https://github.com/nathanmartins/k8s-secret-fmt/releases/download/v1.0.0/k8s-secret-fmt_Linux_x86_64.tar.gz | tar xz
sudo mv k8s-secret-fmt /usr/local/bin/
```

### From Source

If you have Go installed, you can build from source:

```bash
git clone https://github.com/nathanmartins/k8s-secret-fmt.git
cd k8s-secret-fmt
go install
```

### Using Go Install

If you have Go installed, you can install directly using:

```bash
go install github.com/nathanmartins/k8s-secret-fmt@latest
```

## Usage

k8s-secret-fmt reads YAML from standard input and outputs the formatted YAML to standard output. It's particularly useful in combination with SOPS for managing encrypted Kubernetes secrets.

### Basic Usage

```bash
cat secret.yaml | k8s-secret-fmt > formatted-secret.yaml
```

### With SOPS

The most common use case is formatting decrypted SOPS secrets:

```bash
sops -d secrets.enc.yaml | k8s-secret-fmt > secrets.yaml
```

### Example

Input YAML:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
type: Opaque
stringData:
  username: admin
  password: complex"password'with"quotes
  api-key: 1234-5678-abcd
```

Output YAML:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
type: Opaque
stringData:
  username: 'admin'
  password: 'complex"password''with"quotes'
  api-key: '1234-5678-abcd'
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
