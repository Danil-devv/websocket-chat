package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Server struct {
	Port            string `json:"port"`
	WriteBufferSize int    `json:"writeBufferSize"`
	ReadBufferSize  int    `json:"readBufferSize"`
}

type App struct {
	MessagesToLoad int `json:"messagesToLoad"`
}

func LoadApp() (*App, error) {
	size, ok := os.LookupEnv("MESSAGES_TO_LOAD")
	if !ok {
		return nil, errors.New("cannot find 'MESSAGES_TO_LOAD' variable in environment")
	}
	messagesToLoad, err := strconv.Atoi(size)
	if err != nil {
		return nil, fmt.Errorf("%s: variable 'MESSAGES_TO_LOAD' must be integer", err.Error())
	}
	return &App{MessagesToLoad: messagesToLoad}, nil
}

func LoadServer() (*Server, error) {
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
