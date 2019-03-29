package cfg
type RedisConfig struct {
	Url                string `yaml:"url"`
	Auth               string `yaml:"auth"`
	MaxIdle            int    `yaml:"max_idle"`
	IdleTimeOutSeconds int    `yaml:"idle_time_out_seconds"`
}