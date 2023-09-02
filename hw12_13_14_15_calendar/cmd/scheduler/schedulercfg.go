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
	Src       SchedulerSource
	Dest      SchedulerTarget
	Timeout   int
	LoggLevel string
}

func NewSchedulerConfig() *SchedulerCfg {
	return &SchedulerCfg{}
}
