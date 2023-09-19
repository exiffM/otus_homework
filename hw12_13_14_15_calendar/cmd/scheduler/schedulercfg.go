package main

type SchedulerSource struct {
	ConnectionString string
}

type SchedulerTarget struct {
	ConnectionString string
	ExchangeName     string
	Key              string
}

type SchedulerCfg struct {
	Source  SchedulerSource
	Target  SchedulerTarget
	Timeout int
	Logger  string
}

func NewSchedulerConfig() *SchedulerCfg {
	return &SchedulerCfg{}
}
