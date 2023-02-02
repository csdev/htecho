package handlers

type Options struct {
	IncludeAuth bool
}

func DefaultOptions() *Options {
	return &Options{
		IncludeAuth: false,
	}
}
