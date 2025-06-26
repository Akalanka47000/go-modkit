package enums

type options struct {
	lowercase bool
	uppercase bool
}

type optionBuilder func(*options) *options

// Treat enum values as lowercase strings when automatically derived from struct field names
//
// Do not use this option with Uppercase() as it will override it.
func Lowercase() optionBuilder {
	return func(opts *options) *options {
		opts.lowercase = true
		return opts
	}
}

// Treat enum values as uppercase strings when automatically derived from struct field names
//
// Do not use this option with Lowercase() as it will override it.
func Uppercase() optionBuilder {
	return func(opts *options) *options {
		opts.uppercase = true
		return opts
	}
}
