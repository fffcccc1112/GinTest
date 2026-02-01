package config

type Config struct {
	ServerConfig ServerConfig `yaml:"server"`
	AuthConfig   AuthConfig   `yaml:"auth"`
	LoggerConfig LoggerConfig `yaml:"logger"`
	Consul       Consul       `yaml:"Consul"`
	Mysql        Mysql        `yaml:"Mysql"`
	Redis        Redis        `yaml:"Redis"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}
type AuthConfig struct {
	Secret string `yaml:"secret"`
}
type LoggerConfig struct {
	Level    string `yaml:"level"`
	Filename string `yaml:"filename"`
}
type Consul struct {
	Addr string
	Port int64
}
type Mysql struct {
	Addr     string `yaml:"Addr"`
	Port     int    `yaml:"Port"`
	Password string `yaml:"Password"`
	Name     string `yaml:"Name"`
	UserName string `yaml:"UserName"`
}
type Redis struct {
	Addr     string `yaml:"Addr"`
	Passward string `yaml:"Passward"`
	DB       int    `yaml:"DB"`
	PoolSize int    `yaml:"PoolSize"`
}
