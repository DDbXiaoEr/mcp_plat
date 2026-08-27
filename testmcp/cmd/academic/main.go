package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"testmcp/internal/academic"
	"testmcp/pkg/mcp"
)

type Config struct {
	Server struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
		Host    string `yaml:"host"`
		Ports   []int  `yaml:"ports"`
	} `yaml:"server"`
}

func main() {
	configPath := flag.String("config", "configs/academic.yaml", "config file path")
	flag.Parse()

	data, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	srv := mcp.NewServer(cfg.Server.Name, cfg.Server.Version)
	academic.RegisterTools(srv)

	addrs := make([]string, len(cfg.Server.Ports))
	for i, port := range cfg.Server.Ports {
		addrs[i] = fmt.Sprintf("%s:%d", cfg.Server.Host, port)
	}
	log.Fatal(srv.ListenAndServeMulti(addrs))
}
