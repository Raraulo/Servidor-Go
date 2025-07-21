package token

import "os"

var SecretKey = []byte(os.Getenv("123"))
