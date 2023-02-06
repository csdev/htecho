package handlers

type Options struct {
	AccessLog   bool
	IncludeAuth bool
	IncludeIps  bool
}

func DefaultOptions() *Options {
	return &Options{
		AccessLog:   false,
		IncludeAuth: false,
		IncludeIps:  false,
	}
}

func (o *Options) IncludeAll() {
	o.IncludeAuth = true
	o.IncludeIps = true
}
