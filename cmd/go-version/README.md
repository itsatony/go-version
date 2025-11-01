# go-version CLI

Command-line tool for displaying version information from go-version manifests.

## Installation

### Using go install

```bash
go install github.com/itsatony/go-version/cmd/go-version@latest
```

### Building from source

```bash
git clone https://github.com/itsatony/go-version.git
cd go-version/cmd/go-version
go build -o go-version
```

## Usage

### Basic Usage

Display all version information:
```bash
go-version
```

### Options

```
  -manifest string
        Path to versions.yaml manifest file (default: versions.yaml)
  -json
        Output in JSON format
  -compact
        Show compact single-line format
  -schemas
        Show only database schema versions
  -apis
        Show only API versions
  -components
        Show only component versions
  -git
        Show only git information
  -build
        Show only build information
  -help
        Show this help message
```

## Examples

### Show all version information
```bash
go-version
```

Output:
```
Project Information:
  Name:     my-app
  Version:  1.2.3

Git Information:
  Commit:       abc123...
  Tree State:   clean

Build Information:
  Time:        2025-10-11T15:00:00Z
  Go Version:  go1.24.0

Database Schemas:
  postgres_main:  45
  redis_cache:    3
```

### JSON output
```bash
go-version -json
```

Output:
```json
{
  "project": {
    "name": "my-app",
    "version": "1.2.3"
  },
  "git": {
    "commit": "abc123...",
    "tree_state": "clean"
  },
  "schemas": {
    "postgres_main": "45",
    "redis_cache": "3"
  }
}
```

### Compact format
```bash
go-version -compact
```

Output:
```
my-app 1.2.3 (abc123...)
```

### Custom manifest file
```bash
go-version -manifest ./config/versions.yaml
```

### Show only schemas
```bash
go-version -schemas
```

Output:
```
Database Schemas:
  postgres_main:  45
  redis_cache:    3
```

### Show only git information
```bash
go-version -git
```

Output:
```
Git Information:
  Commit:       abc123...
  Tree State:   clean
  Commit Time:  2025-10-11T14:55:23Z
```

### Show only APIs
```bash
go-version -apis
```

Output:
```
API Versions:
  rest_v1:  1.15.0
  grpc:     1.2.0
```

## Use Cases

### CI/CD Pipelines

Check version before deployment:
```bash
VERSION=$(go-version -compact)
echo "Deploying version: $VERSION"
```

Get version as JSON for processing:
```bash
go-version -json | jq -r '.project.version'
```

### Shell Scripts

Extract specific information:
```bash
# Get project version
VERSION=$(go-version -json | jq -r '.project.version')

# Get git commit
COMMIT=$(go-version -json | jq -r '.git.commit')

# Check schema version
SCHEMA_VERSION=$(go-version -json | jq -r '.schemas.postgres_main')
```

### Docker Builds

Include version info in Docker labels:
```dockerfile
RUN go-version -json > /app/version.json

LABEL version="$(go-version -compact)"
LABEL git.commit="$(go-version -json | jq -r '.git.commit')"
```

### Debugging

Quick version check:
```bash
# Check what version is deployed
ssh production-server 'cd /app && ./go-version'

# Compare local vs deployed version
diff <(go-version -json) <(ssh prod './go-version -json')
```

### Makefile Integration

```makefile
.PHONY: version
version:
	@go-version -compact

.PHONY: version-json
version-json:
	@go-version -json

.PHONY: check-schema
check-schema:
	@go-version -schemas
```

## Exit Codes

- `0` - Success
- `1` - Error (manifest not found, invalid format, etc.)

## Environment

The CLI tool respects the same environment variables as the go-version library:
- Searches for `versions.yaml` in the current directory by default
- Use `-manifest` flag to specify a custom path
- Works with any valid go-version manifest file

## Building with Version Info

Build the CLI with embedded version information:

```bash
go build -ldflags="\
  -X github.com/itsatony/go-version.GitCommit=$(git rev-parse HEAD) \
  -X github.com/itsatony/go-version.GitTag=$(git describe --tags --always) \
  -X github.com/itsatony/go-version.BuildTime=$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
  -o go-version
```

## License

MIT License - see [LICENSE](../../LICENSE) file for details.
