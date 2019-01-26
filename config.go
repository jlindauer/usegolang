package main

import (
  "fmt"
  "encoding/json"
  "os"
)

type PostgresConfig struct {
  Host     string `json:"host"`
  Port     int    `json:"port"`
  User     string `json:"user"`
  Password string `json:"password"`
  Name     string `json:"name"`
}

func (c PostgresConfig) Dialect() string {
  return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
  // Provide two potential connection info strings based on whether a password is present
  if c.Password == "" {
    return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
      c.Host, c.Port, c.User, c.Name)
  }
  return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    c.Host, c.Port, c.User, c.Password, c.Name)
}

func DefaultPostgresConfig() PostgresConfig {
  return PostgresConfig{
    Host:     "localhost",
    Port:     5432,
    User:     "postgres",
    Password: "",
    Name:     "usegolang_dev",
  }
}

type Config struct {
  Port     int            `json:"port"`
  Env      string         `json:"env"`
  Pepper   string         `json:"pepper"`
  HMACKey  string         `json:"hmac_key"`
  Database PostgresConfig `json:"database"`
}

func (c Config) IsProd() bool {
  return c.Env == "prod"
}

func DefaultConfig() Config {
  return Config{
    Port:     3000,
    Env:      "dev",
    Pepper:   "secret-random-string",
    HMACKey:  "secret-hmac-key",
    Database: DefaultPostgresConfig(),
  }
}

func LoadConfig(configReq bool) Config {
  // Open the config file
  f, err := os.Open(".config")
  if err != nil {
    if configReq {
      panic(err)
    }
    fmt.Println("Using the default config...")
    return DefaultConfig()
  }
  // If successful create a Config variable to load into
  var c Config
  dec := json.NewDecoder(f)
  err = dec.Decode(&c)
  if err != nil {
    panic(err)
  }
  fmt.Println("Successfully loaded .config")
  return c
}
