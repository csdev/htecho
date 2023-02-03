package handlers

type Options struct {
	IncludeAuth bool
	IncludeIps  bool
}

func DefaultOptions() *Options {
	return &Options{
		IncludeAuth: false,
		IncludeIps:  false,
	}
}

func (o *Options) IncludeAll() {
	o.IncludeAuth = true
	o.IncludeIps = true
}
