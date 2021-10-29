package utils

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid"
)

const (
	base       = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	baseLength = len(base)
)

func NewID() string {
	return shortuuid.New()
}

func NewUUID() string {
	return uuid.New().String()
}

func Md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func RandString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = base[rand.Intn(baseLength)]
	}
	randString := string(b)
	return randString
}

func RandDigists(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = byte('0' + rand.Intn(10))
	}
	randDigists := string(b)
	return randDigists
}
