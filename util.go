package main

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"time"
)

func TokenHex(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	h := sha1.New()
	h.Write(b)
	sum := h.Sum(nil)
	return fmt.Sprintf("%x", sum)
}

func RandBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}
