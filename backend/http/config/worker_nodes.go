package config

import (
	"os"
	"strings"
)

type WorkerNodesConfig struct {
	Addresses []string
}

func LoadWorkerNodesConfig() *WorkerNodesConfig {
	return &WorkerNodesConfig{
		Addresses: strings.FieldsFunc(os.Getenv("WORKER_NODES"), func(c rune) bool {
			return c == ','
		}),
	}
}
