package scheduler

type Source struct {
	ConnectionString string
}

type Target struct {
	ConnectionString string
	ExchangeName     string
	Key              string
}

type Cfg struct {
	Source  Source
	Target  Target
	Timeout int
	Logger  string
}

func NewSchedulerConfig() *Cfg {
	return &Cfg{}
}
