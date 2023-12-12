package client

type OptFunc func(*opts)

type opts struct {
	// regex indicates whether the key should be treated as a regex.
	//
	// Applicable for: GET, PUT, DEL, WATCH
	regex bool
	// keysOnly specifies whether only keys should be returned without
	// values.
	//
	// Applicable for: GET, WATCH
	keysOnly bool
}

// WithRegex configures the operation to treat the key as a regex.
//
// Applicable for: GET, PUT, DEL, WATCH
func WithRegex() OptFunc {
	return func(o *opts) {
		o.regex = true
	}
}

// WithKeysOnly configures the operation to only return keys without
// values.
//
// Applicable for: GET, WATCH
func WithKeysOnly() OptFunc {
	return func(o *opts) {
		o.keysOnly = true
	}
}

func defaultOpts() opts {
	return opts{
		regex:    false,
		keysOnly: false,
	}
}
