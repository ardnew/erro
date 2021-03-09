package cmd

// Mode represents the configured mode of operation.
type Mode int

// Constants values of enumerated type Mode.
const (
	Literal Mode = iota
	Escaped
	Formats
)

// Config contains the configuration options and all arguments for a single
// invocation of the erro command.
type Config struct {
	ShowVersion bool
	ShowChanges bool
	AddNewline  bool
	Mode        Mode
	Format      string
	Args        []string
}
