package Cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	// "encoding/hex"
	// "fmt"
	"io"
	"os"
)
type CTR struct {
	SecretKey string
}

// inspiration: https://golang.org/src/crypto/cipher/example_test.go

// func main() {
// 	c := CTR{
// 		"6368616e676520746869732070617373",
// 	}

// 	bytesToBeEncrypted := []byte("hello there")

// 	encryptedBytes := c.encrypt(bytesToBeEncrypted)
// 	decryptedBytes := c.decrypt(encryptedBytes)

// 	fmt.Println(decryptedBytes)
// }

func (c *CTR) EncryptToFile(name, plaintext string) {
	cipherBytes := c.Encrypt([]byte(plaintext))
	err := os.WriteFile(name, cipherBytes, 0644)
	check(err)
}

func (c *CTR) DecryptFromFile(name string) string {
	cipherBytes, err := os.ReadFile(name)
	check(err)
	decryptedBytes := c.Decrypt(cipherBytes)
	return string( decryptedBytes[:] )
}

func (c *CTR) Encrypt(inputBytes []byte) []byte {
	// key, _ := hex.DecodeString(c.SecretKey)
	key := []byte(c.SecretKey)
	
	block, err := aes.NewCipher(key)
	check(err)

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the byte-array
	outputBytes := make([]byte, aes.BlockSize+len(inputBytes))

	iv := outputBytes[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	check(err)

	stream := cipher.NewCTR(block, iv)
	// fmt.Println(len(outputBytes[aes.BlockSize:]))
	// fmt.Println(len(inputBytes))

	stream.XORKeyStream(outputBytes[aes.BlockSize:], inputBytes)
	// fmt.Println(len(outputBytes[aes.BlockSize:]))

	// fmt.Printf("encrypt: \n'%s'\n'%s'\n\n", string(inputBytes[:]), string(outputBytes[:]))

	return outputBytes
}
func (c *CTR) Decrypt(inputBytes []byte) []byte {
	// key, _ := hex.DecodeString(c.SecretKey)
	key := []byte(c.SecretKey)

	block, err := aes.NewCipher(key)
	check(err)

	content := inputBytes[aes.BlockSize:]
	iv := inputBytes[:aes.BlockSize]

	outputBytes := make([]byte, len(content))

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(outputBytes, content)

	// fmt.Printf("decrypt: \n'%s'\n'%s'\n\n", string(inputBytes[:]), string(outputBytes[:]))
	// fmt.Println(len(outputBytes))

	return outputBytes
}
func GenerateNewRndmIV(size int) string {
	iv := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, iv)
	check(err)
	return string(iv)
}

func check(e error) {
	if e != nil {panic(e)}
}
