package ants

// contains definitions not defined in the original ants source, but required
// by pants.

type (
	// IDGenerator is a sequential unique id generator interface
	IDGenerator interface {
		Generate() string
	}
)
