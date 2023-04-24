# Crda CLI 1.5

This project is still in development mode.
Staging version (Pre-release) can be found [here](https://github.com/RHEcosystemAppEng/crda-cli/releases/tag/staging).
Download you binary based on your OS.

```shell
$ crda analyse /path/to/maven/project/pom.xml

Full Report:  file:///tmp/crda/stack-analysis-maven-1682328584.html
```

> Currently, only Java's Maven ecosystem is implemented.

```shell
$ crda help

Use this tool for CodeReady Dependency Analytics reports

Usage:
  crda [command]

Available Commands:
  analyse     Preform dependency analysis report
  auth        Link crda user with snyk
  completion  Generate a completions script
  config      Manage crda config
  help        Help about any command
  version     Get binary version

Flags:
  -m, --client string   The invoking client for telemetry (default "terminal")
  -d, --debug           Set DEBUG log level
  -h, --help            help for crda
  -c, --no-color        Toggle colors in output.

Use "crda [command] --help" for more information about a command.
```
