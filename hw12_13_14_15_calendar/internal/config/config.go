package config

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type HTTPConf struct {
	Host string
	Port string
}

type LoggerConf struct {
	Level string
}

type StorageConf struct {
	DSN  string
	Type string
}

type Config struct {
	HTTP    HTTPConf
	Logger  LoggerConf
	Storage StorageConf
}

func NewConfig() *Config {
	return &Config{}
}
