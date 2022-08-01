# CLI
## Installation
Install `tm` by running following command:
```shell
go install github.com/yule-l/tm/cmd/tm@latest
```

and ensuring that `$GOPATH/bin` is added to your `$PATH`.

## Usage
Just run 
```shell
tm do --input tasks.txt --script ./do.sh --parallel --force
```

And task manager will get your tasks and put then one by one into `./do.sh` as first argument.
You can read more about flags in the help `tm do --help`.

## Tutorial

### First step, prepare your tasks file

See [tasks.txt](tasks.txt) for an example.

```text
first
second
third
fourth
fifth
sixth
seventh
eighth
ninth
tenth
eleventh
twelfth
thirteenth
fourteenth
fifteenth
```

One line will be used as a first argument for the script.

### Second step, prepare your script

See [do.sh](do.sh) for an example.

```bash
#!/usr/bin/env bash

printf "processing task %s" "$1"
```

### Third step, run tm

By default `tm do` will run your tasks from in parallel, it will use file tasks.txt as input and will use `do.sh` as script.

```shell
tm do
```

expected output:

```text  
processing task fifth
processing task second
processing task fourth
processing task third
processing task first
processing task sixth
processing task eighth
processing task ninth
processing task tenth
processing task seventh
processing task eleventh
processing task twelfth
processing task fourteenth
processing task thirteenth
processing task fifteenth
```