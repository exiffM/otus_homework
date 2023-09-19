package main

type SenderSource struct {
	ConnectionString string
	QueueName        string
	ExchangeName     string
	Key              string
	Tag              string
}

type SenderCfg struct {
	Source    SenderSource
	LoggLevel string
}

func NewSenderConfig() *SenderCfg {
	return &SenderCfg{}
}
