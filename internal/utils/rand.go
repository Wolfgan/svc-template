package utils

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
	"unsafe"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var source = rand.NewSource(time.Now().UnixNano())

// RandString генерирует случайную строку заданной длинны.
func RandString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, source.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = source.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func SHA256(bytes []byte) string {
	sh := sha256.New()
	sh.Write(bytes)
	return hex.EncodeToString(sh.Sum(nil))
}

func SHA1(bytes []byte) string {
	sh := sha1.New()
	sh.Write(bytes)
	return hex.EncodeToString(sh.Sum(nil))
}
