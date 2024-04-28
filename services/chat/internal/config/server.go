package config

import (
	"chat/internal/adapters/websocket"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Server struct {
	Port            string
	WriteBufferSize int
	ReadBufferSize  int
}

func getServerConfig() (*websocket.Config, error) {
	cfg, err := loadEnvServerConfig()
	if err != nil {
		return nil, err
	}
	return &websocket.Config{
		Port:            cfg.Port,
		WriteBufferSize: cfg.WriteBufferSize,
		ReadBufferSize:  cfg.ReadBufferSize,
	}, nil
}

func loadEnvServerConfig() (*Server, error) {
	port, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		return nil, errors.New("cannot find 'PORT' variable in environment")
	}

	size, ok := os.LookupEnv("WRITE_BUFFER_SIZE")
	if !ok {
		return nil, errors.New("cannot find 'WRITE_BUFFER_SIZE' variable in environment")
	}
	writeBufferSize, err := strconv.Atoi(size)
	if err != nil {
		return nil, fmt.Errorf("%s: variable 'WRITE_BUFFER_SIZE' must be integer", err.Error())
	}

	size, ok = os.LookupEnv("READ_BUFFER_SIZE")
	if !ok {
		return nil, errors.New("cannot find 'READ_BUFFER_SIZE' variable in environment")
	}

	readBufferSize, err := strconv.Atoi(size)
	if err != nil {
		return nil, fmt.Errorf("%s: variable 'READ_BUFFER_SIZE' must be integer", err.Error())
	}

	return &Server{
		Port:            port,
		WriteBufferSize: writeBufferSize,
		ReadBufferSize:  readBufferSize,
	}, nil
}
