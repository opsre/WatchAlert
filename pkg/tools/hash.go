package tools

import (
	"crypto/md5"
	"fmt"
)

const (
	offset64 = 14695981039346656037
	prime64  = 1099511628211
)

// HashNew initializes a new fnv64a hash value.
func HashNew() uint64 {
	return offset64
}

// HashAdd adds a string to a fnv64a hash value, returning the updated hash.
func HashAdd(h uint64, s string) uint64 {
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= prime64
	}
	return h
}

// HashAddByte adds a byte to a fnv64a hash value, returning the updated hash.
func HashAddByte(h uint64, b byte) uint64 {
	h ^= uint64(b)
	h *= prime64
	return h
}

func Md5Hash(str []byte) string {

	hash := md5.Sum(str)
	hashStr := fmt.Sprintf("%x", hash)
	return hashStr

}

func WithKVCalculateHash(key, value string) string {
	return Md5Hash([]byte(key + ":" + value))
}
