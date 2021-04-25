package middleware

import (
	"log"
	"os"
)

var (
	middlwareLogger = log.New(os.Stdout, "MIDDLEWARE:", log.LstdFlags)
)
