package handler

import (
	"crypto/md5"
	"fmt"
	"time"
)

func generateLink() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String())))[:15]
}
