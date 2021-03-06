package partition

import (
	"github.com/google/uuid"
)

// Options is the functional options struct.
type Options struct {
	Type uuid.UUID
	Name string
}

// Option is the functional option func.
type Option func(*Options)

// WithPartitionType sets the partition type.
func WithPartitionType(o [16]byte) Option {
	return func(args *Options) {
		// TODO: An Option should return an error.
		// nolint: errcheck
		guuid, _ := uuid.FromBytes(o[:])
		args.Type = guuid
	}
}

// WithPartitionName sets the partition name.
func WithPartitionName(o string) Option {
	return func(args *Options) {
		args.Name = o
	}
}

// NewDefaultOptions initializes a Options struct with default values.
func NewDefaultOptions(setters ...interface{}) *Options {
	// Default to data type "af3dc60f-8384-7247-8e79-3d69d8477de4"
	// TODO: An Option should return an error.
	// nolint: errcheck
	guuid, _ := uuid.FromBytes([]byte{0Xaf, 0X3d, 0Xc6, 0X0f, 0X83, 0X84, 0X72, 0X47, 0X8e, 0X79, 0X3d, 0X69, 0Xd8, 0X47, 0X7d, 0Xe4})

	opts := &Options{
		Type: guuid,
		Name: "",
	}

	for _, setter := range setters {
		if s, ok := setter.(Option); ok {
			s(opts)
		}
	}

	return opts
}
