package config

import (
	"os"
	"strconv"
)
import "gopkg.in/yaml.v3"

var GConfig *Config

// 加载配置文件
func Load() *Config {
	file, err := os.Open("D:\\project\\microProject\\ginTest\\config\\config.yaml")
	if err != nil {
		panic("读取配置文件失败" + err.Error())
	}
	defer file.Close()
	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		panic("读取配置文件失败" + err.Error())
	}
	//环境变量覆盖
	if port := os.Getenv("SERVER_PORT"); port != "" {
		portInt, err := strconv.Atoi(port)
		if err != nil {
			print("port转int失败")
		}
		cfg.ServerConfig.Port = portInt
	}
	if secret := os.Getenv("AUTH_SECRET"); secret != "" {
		cfg.AuthConfig.Secret = secret
	}
	GConfig = &cfg
	return &cfg
}
