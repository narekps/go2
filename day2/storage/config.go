package storage

type Config struct {
	Dsn string
}

func NewConfig() *Config {
	return &Config{
		Dsn: "./database.db",
	}
}
