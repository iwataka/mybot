package models

import (
	"time"
)

type Context interface {
	String(key string) string
	Bool(key string) bool
	Int(key string) int
	Duration(key string) time.Duration
}
