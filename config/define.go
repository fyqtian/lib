package config

import "time"

type Configer interface {
	GetString(string) string
	GetInt(string) int
	GetBool(string) bool
	GetDuration(string) time.Duration
}
