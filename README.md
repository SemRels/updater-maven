# updater-maven

Updates the project version in a Maven `pom.xml` file.

This plugin is distributed as the standalone Go binary `semrel-plugin-updater-maven`. Semrel executes the binary as a subprocess, provides plugin configuration through `SEMREL_PLUGIN_*` environment variables, provides release context through `SEMREL_*` environment variables, reads standard output, and treats exit code `0` as success and any non-zero exit code as failure. Install the binary in `~/.semrel/plugins/` or anywhere on your `$PATH`.

## Installation

```bash
go install github.com/SemRels/updater-maven/cmd/plugin@latest
```

## Configuration

```yaml
plugins:
  - name: updater-maven
    path: ~/.semrel/plugins/semrel-plugin-updater-maven
    env:
      SEMREL_PLUGIN_FILE: "pom.xml"
```

## `SEMREL_PLUGIN_*` variables

| Name | Required | Description | Default |
| --- | --- | --- | --- |
| `SEMREL_PLUGIN_FILE` | Optional | Path to the Maven POM file to update. | pom.xml |

## `SEMREL_*` release context used

| Variable | Description |
| --- | --- |
| `SEMREL_VERSION` | Resolved release version for the current run. |
| `SEMREL_NEXT_VERSION` | Next version computed by semrel for the release. |
| `SEMREL_DRY_RUN` | Whether semrel is running in dry-run mode. |

## Example behavior

The plugin rewrites the Maven project version to the next release version and reports the updated file.

## License

Apache-2.0
