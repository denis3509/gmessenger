package config
 
type Config struct {
	Port int
	Debug bool
	DSN string
}
var Default = Config {
	Port: 3000,
	Debug : false,
	DSN: "postgres://localhost:5432/golang_messenger?sslmode=disable&user=denis&password=localpass",
}

func GetConfig() Config{ 
	return Default
}

func LoadFromDotEnv(path string) *Config {
	return nil
}
