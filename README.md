# ENVP

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/sunggun-yu/envp)
![GitHub all releases](https://img.shields.io/github/downloads/sunggun-yu/envp/total)
![test](https://github.com/sunggun-yu/envp/actions/workflows/test.yaml/badge.svg)
![release](https://github.com/sunggun-yu/envp/actions/workflows/release.yaml/badge.svg)
[![CodeQL](https://github.com/sunggun-yu/envp/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/sunggun-yu/envp/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sunggun-yu/envp)](https://goreportcard.com/report/github.com/sunggun-yu/envp)
[![codecov](https://codecov.io/gh/sunggun-yu/envp/branch/main/graph/badge.svg?token=3V5SJ002BS)](https://codecov.io/gh/sunggun-yu/envp)

`ENVP` is a shell wrapper command-line tool for macOS and Linux that enables you to run commands or shells with specific environment variable configurations based on profiles. It also allows you to run scripts to set up these profiles. With ENVP, you can easily switch between different environment configurations, even simultaneously in different terminal windows. This tool provides a convenient way to manage and control your environment variables for various development and testing scenarios.

![envp intro](docs/assets/envp-intro.gif)

## Installation

brew:

```bash
brew install sunggun-yu/tap/envp
```

go install:

```bash
go install github.com/sunggun-yu/envp@<version>
```

## Quick Start

### Config file

Location of config file is `~/.config/envp/config.yaml`. 

> please create folder and config file if it is not created.
>
> ```bash
> mkdir -p ~/.config/envp
> vim ~/.config/envp/config.yaml
> ```

config file example:

```yaml
default: ""
profiles:
  my-profile:
    desc: profile description
    env:
      - name: HTTPS_PROXY
        value: http://some-proxy:3128
      - name: MY_PASSWORD
        value: $(cat ~/.config/my-password-file)
    init-script:
      - run: |
          echo "this is init script 1"
      - run: |
          echo "this is init script 2"
```

### Start new shell with profile

You can create new shell session with injected environment variable from your profile.

```bash
# start new shell session with specific profile
envp start my-profile
```
