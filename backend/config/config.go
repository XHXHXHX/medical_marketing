package config

type Config struct {
	Server struct {
		Addr string `yaml:"addr"`
	} `yaml:"server"`
	Mysql struct {
		Host string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName string `yaml:"db_name"`
	} `yaml:"mysql"`
	Redis struct {
		Host   string `yaml:"host"`
		Passwd string `yaml:"password"`
		DB     int    `yaml:"db"`
	} `yaml:"redis"`
}