# tm – task manager
[![GoDoc](https://pkg.go.dev/badge/github.com/yule-l/tm)](https://pkg.go.dev/github.com/yule-l/tm)
[![Go](https://github.com/yule-l/tm/actions/workflows/go.yml/badge.svg)](https://github.com/yule-l/tm/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/yule-l/tm/branch/master/graph/badge.svg?token=ZRL9IO6JNC)](https://codecov.io/gh/yule-l/tm)
[![Go Report Card](https://goreportcard.com/badge/github.com/yule-l/tm)](https://goreportcard.com/report/github.com/yule-l/tm)

The simplest task manager in Go.

## Overview

You can use `tm` as a library or as a cli manager, ones providing a simple mechanism for tasks control.

### Features
* Command Line Interface – use cli to do your tasks
* Read tasks line by line from file
* Mark tasks as completed
* Mark tasks as not completed and try complete this N times in the future
* Preconfigured structs for parallel and ordered tasks executions

## CLI Installing

Install `tm` by running:

```shell
go install github.com/yule-l/tm/cmd/tm@latest
```

## Documentation
- [CLI][]
- [Library][]

[CLI]: ./docs/cli/README.md
[Library]: ./docs/lib/README.md

## License

This library is licensed under either of

* [Apache License, Version 2.0](LICENSE-APACHE)
* [MIT license](LICENSE-MIT)

at your option.
