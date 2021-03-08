[docimg]:https://godoc.org/github.com/ardnew/erro?status.svg
[docurl]:https://godoc.org/github.com/ardnew/erro
[repimg]:https://goreportcard.com/badge/github.com/ardnew/erro
[repurl]:https://goreportcard.com/report/github.com/ardnew/erro

# erro
#### Echo text to stderr

[![GoDoc][docimg]][docurl] [![Go Report Card][repimg]][repurl]

## Usage

Use `erro` just like you would `echo`:

```
$ erro "Hello darkness, my old friend"
Hello darkness, my old friend
```

As a formatting convenience, you can also perform argument substitution with the `-f` flag as follows:

```
$ erro -f "Hello {1}, my old {0}" darkness friend
Hello friend, my old darkness
```

Use the `-h` flag for usage summary:

```sh
usage:
  erro [options] [args ...]

options:
  -v
        Display version information
  -V
        Display change history
  -n
        Do not output a trailing newline
  -e
        Enable interpretation of backslash escapes
  -E
        Disable interpretation of backslash escapes (default true)
  -f format
        Format output string according to format, where "{N}" represents argument N
```

## Installation

Use the builtin Go package manager:

```sh
go get -v github.com/ardnew/erro
```
