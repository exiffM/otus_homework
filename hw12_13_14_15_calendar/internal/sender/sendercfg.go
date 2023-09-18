package sender

type Source struct {
	ConnectionString string
	QueueName        string
	ExchangeName     string
	Key              string
	Tag              string
}

type Cfg struct {
	Source    Source
	LoggLevel string
}

func NewSenderConfig() *Cfg {
	return &Cfg{}
}
