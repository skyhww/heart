package cfg

type DatabaseConfig struct {
	Driver                 string `yaml:"driver"`
	Dsn                    string `yaml:"dsn"`
	ConnMaxLifetimeSeconds int    `yaml:"conn_max_lifetime_seconds"`
	MaxIdleConns           int    `yaml:"max_idle_conns"`
	MaxOpenConns           int    `yaml:"max_open_conns"`
}