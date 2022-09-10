package config

import "fmt"

 

type DBConfig struct {
	Name      string
	DriverName string
	Host string
	Port      int
	User      string
	Password  string
}

type Config struct {
	DB DBConfig
	Port  int
	SocketPort int
	Debug bool
	DSN   string
}

var Default = Config{
	Port:  8000,
 
	Debug: false,
	DSN:   "postgres://localhost:5432/golang_messenger?sslmode=disable&user=denis&password=localpass",
	DB: DBConfig {
		Name:      "golang_messenger",
		DriverName: "postgres",
		Port:     5432,
		User:     "denis",
		Password: "localpass",
	},
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("%s://%s:%d/%s?sslmode=disable&user=%s&password=%s",
		c.DriverName,
		c.Host,
		c.Port,
		c.Name,
		c.User,
		c.Password,
	)
}



func GetConfig() Config {
	return Default
}

func LoadFromDotEnv(path string) *Config {
	return nil
}
