package config

const (
	RedisEnable = true
	RedisMode   = "single" // "single" or "cluster"
	RedisDB     = 0
	RedisPass   = ""
)

// For cluster
var RedisAddrs = []string{"localhost:6379"}
