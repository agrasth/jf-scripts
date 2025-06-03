# JFrog CLI Automation Scripts

This repository contains automation scripts for maintaining JFrog CLI components.

## Scripts

### `update_jfrog_versions.go`

Automatically updates all JFrog repository versions in `jfrog-cli-artifactory/go.mod` to their latest GitHub releases.

**Usage:**
```bash
git clone https://github.com/your-username/jf-scripts.git
cd jf-scripts
go run update_jfrog_versions.go
```

**Features:**
- ✅ Fetches latest versions from GitHub API
- ✅ Updates both `require` and `replace` directives
- ✅ Shows current version → new version comparisons
- ✅ Automatically runs `go mod tidy`
- ✅ Skips already up-to-date packages
- ✅ Proper error handling and user feedback

**Requirements:**
- `jfrog-cli-artifactory` directory must be in a sibling directory
- Optional: Set `GITHUB_TOKEN` environment variable to avoid rate limits

**Repositories Updated:**
- `github.com/jfrog/build-info-go`
- `github.com/jfrog/froggit-go`
- `github.com/jfrog/gofrog`
- `github.com/jfrog/jfrog-cli-core/v2`
- `github.com/jfrog/jfrog-client-go`

## Directory Structure

```
workspace/
├── jf-scripts/                  (this repository)
│   ├── update_jfrog_versions.go
│   └── README.md
└── forked/                      (or any parent directory)
    └── jfrog-cli-artifactory/   (target directory)
        └── go.mod               (file to update)
``` 