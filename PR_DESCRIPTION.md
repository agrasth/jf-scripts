# Add Automation Script for JFrog Repository Version Updates

## Description

This PR introduces an automation script that simplifies the process of updating JFrog repository versions in `jfrog-cli-artifactory`. Previously, developers had to manually check and update each JFrog dependency version, which was time-consuming and error-prone.

## Changes

### New Files
- **`automation-scripts/update_jfrog_versions.go`** - Main automation script
- **`automation-scripts/go.mod`** - Go module configuration for the script
- **`automation-scripts/README.md`** - Documentation and usage instructions

### Features Added
- **Automated Version Fetching**: Queries GitHub API to get the latest release for each JFrog repository
- **Smart Updates**: Only updates packages that have newer versions available
- **Comprehensive Coverage**: Updates both `require` and `replace` directives in go.mod
- **Automatic Cleanup**: Runs `go mod tidy` automatically after updates
- **User-Friendly Output**: Clear progress indicators and version change summaries
- **Error Handling**: Graceful handling of API failures and missing dependencies

## Usage

```bash
cd automation-scripts
go run update_jfrog_versions.go
```

## Benefits

1. **Time Saving**: Reduces manual version checking from minutes to seconds
2. **Accuracy**: Eliminates human errors in version updates
3. **Consistency**: Ensures all JFrog dependencies are updated uniformly
4. **Automation Ready**: Can be integrated into CI/CD pipelines
5. **Developer Experience**: Simple one-command operation with clear feedback

## Repositories Managed

The script automatically updates these JFrog repositories:
- `github.com/jfrog/build-info-go`
- `github.com/jfrog/froggit-go`
- `github.com/jfrog/gofrog`
- `github.com/jfrog/jfrog-cli-core/v2`
- `github.com/jfrog/jfrog-client-go`

## Example Output

```
🚀 JFrog CLI Artifactory Version Updater
==========================================
📁 Target: /path/to/jfrog-cli-artifactory/go.mod

🔍 Checking github.com/jfrog/build-info-go...
✅ Updated github.com/jfrog/build-info-go: v1.10.11 → v1.10.12
🔍 Checking github.com/jfrog/froggit-go...
✨ github.com/jfrog/froggit-go is already up to date (v1.17.0)
...
🎉 Successfully updated 3 JFrog repositories!
🧹 Running 'go mod tidy'...
✅ Successfully ran 'go mod tidy'
🎊 All done! JFrog repository versions updated and dependencies cleaned up.
```

## Testing

The script has been tested with:
- Various network conditions (with/without GitHub token)
- Different go.mod states (up-to-date vs outdated packages)
- Error scenarios (missing directories, malformed go.mod files)

## Future Enhancements

- Support for updating other JFrog CLI components
- Integration with GitHub Actions for automated PRs
- Configuration file support for custom repository lists
- Rollback functionality for version downgrades 