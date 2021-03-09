[docimg]:https://godoc.org/github.com/ardnew/erro?status.svg
[docurl]:https://godoc.org/github.com/ardnew/erro
[repimg]:https://goreportcard.com/badge/github.com/ardnew/erro
[repurl]:https://goreportcard.com/report/github.com/ardnew/erro

# erro
#### Echo text to stderr

[![GoDoc][docimg]][docurl] [![Go Report Card][repimg]][repurl]

## Usage

`erro` is a functionally identical clone to `echo` from [GNU coreutils](https://github.com/coreutils/coreutils), except that it prints to the standard error stream (`stderr`) instead of standard output (`stdout`):

```
$ erro "Hello darkness, my old friend"
Hello darkness, my old friend
```

All of the same flags have the same meaning. But as a formatting convenience, you can also perform argument substitution with the `-f` flag as follows:

```
$ erro -f "Hello {1}, my old {0}" darkness friend
Hello friend, my old darkness
```

Use the `-h` flag for usage summary:

```sh
usage:
  erro [options] [args ...]

options:
  -v       Display version information
  -V       Display change history
  -n       Do not output a trailing newline
  -e       Enable interpretation of backslash escapes
  -E       Disable interpretation of backslash escapes (default true)
             Accepted for compatibility with echo, but this flag is ignored.
  -f fmt   Output string according to fmt (implies -e)
             See formatting section below for details.

formatting:
  The fmt argument given to flag -f may contain special placeholder symbols
  of the form "{N}", meaning the N'th command-line argument. The indexing does
  not include flags or their arguments; it is strictly the N'th argument that
  would be printed to the standard error stream if no flags were given at all.

  Additionally, the -f flag enables interpretation of backslash sequences as
  if the -e flag was given.

  The following backslash escape sequences are recognized:

    \\     backslash                        \x5C
    \a     alert                            \x07  BEL
    \b     backspace                        \x08  BS
    \c     produce no further output
    \e     escape                           \x1B  ESC
    \f     form feed                        \x0C  FF
    \n     new line                         \x0A  LF
    \r     carriage return                  \x0D  CR
    \t     horizontal tab                   \x09  TAB
    \v     vertical tab                     \x0B  VT
    \0NNN  byte with octal value NNN        1 to 3 digits
    \xHH   byte with hexadecimal value HH   1 to 2 digits
```

## Installation

Use the builtin Go package manager:

```sh
go get -v github.com/ardnew/erro
```
