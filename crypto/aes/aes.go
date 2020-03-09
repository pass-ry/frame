// AES-256
// Zero-Padding

/*
	WARNING!!!
	Not Finish Yet !!!
	Do Not Use It !!!
*/
package aes

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"

	"github.com/pkg/errors"
	"gitlab.ifchange.com/data/cordwood/util/random"
)

// GenerateKey generate a key of AES-256
// key's size is 32 bytes
func GenerateKey() string {
	return random.RandStr(32)
}

func Encrypt(key string, plaintext string) (ciphertext string, err error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", errors.Wrapf(err, "aes.NewCipher(%s)", key)
	}
	plaintextBytes := []byte(plaintext)
	size := len(plaintextBytes)
	if size%16 != 0 {
		size += 16 - size%16
	}
	paddingPlaintextBytes := make([]byte, size)
	copy(paddingPlaintextBytes, plaintextBytes)
	ciphertextBytes := make([]byte, size)
	c.Encrypt(ciphertextBytes, paddingPlaintextBytes)
	ciphertext = hex.EncodeToString(ciphertextBytes)
	return
}

func Decrypt(key string, ciphertext string) (plaintext string, err error) {
	ciphertextBytes, _ := hex.DecodeString(ciphertext)
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", errors.Wrapf(err, "aes.NewCipher(%s)", key)
	}
	plaintextBytes := make([]byte, len(ciphertextBytes))
	c.Decrypt(plaintextBytes, ciphertextBytes)
	plaintext = string(plaintextBytes[:])
	return
}

func padding(src []byte, blocksize int) []byte {
	padnum := blocksize - len(src)%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	return append(src, pad...)
}

func unpadding(src []byte) []byte {
	n := len(src)
	unpadnum := int(src[n-1])
	return src[:n-unpadnum]
}
