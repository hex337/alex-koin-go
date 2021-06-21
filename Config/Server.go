package Config

import (
	"fmt"
	"os"
)

type ServerConfig struct {
	Host string
	Port string
}

func BuildServerConfig() *ServerConfig {
	port, portProvided := os.LookupEnv("PORT")
	if !portProvided {
		port = "3000"
	}

	serverConfig := ServerConfig{
		Host: "localhost",
		Port: port,
	}
	return &serverConfig
}

func ServerURL(serverConfig *ServerConfig) string {
	return fmt.Sprintf(
		"%s:%s",
		serverConfig.Host,
		serverConfig.Port,
	)
}
