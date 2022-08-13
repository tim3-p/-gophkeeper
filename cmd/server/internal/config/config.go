package config

// Config contains client config parameters set in the config file
type Config struct {
	StoreFile  string `json:"store_file" env:"STORE_FILE" envDefault:"server_storage.db"`
	ServerKey  string `json:"server_key" env:"SERVER_KEY" envDefault:"keys/server.key"`
	ServerCRT  string `json:"server_crt" env:"SERVER_CRT" envDefault:"keys/server.crt"`
	ServerPort int    `json:"server_port" env:"SERVER_PORT" envDefault:"8080"`
}

// Cfg holds global parameters from config file
var Cfg Config
