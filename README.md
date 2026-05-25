# updater-maven

`updater-maven` is a SemRel plugin that updates Maven project versions and deploys artifacts.

## Configuration

Environment variables:

- `MAVEN_SETTINGS` (optional path to `settings.xml`)

## Behavior

The plugin runs:

1. `mvn -B versions:set -DnewVersion=<version> -DgenerateBackupPoms=false`
2. `mvn -B deploy -DskipTests`

If `MAVEN_SETTINGS` is set, the plugin passes `--settings <path>` to both commands.

## Development

```bash
go mod tidy
go build ./...
go test ./...
```
