package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)
type CBC struct {
	Iv string
}

func main() {
	m := make(map[int]CBC)
	m[2] = CBC{Iv: "test2"}
	m[4] = CBC{Iv: "test4"}
	fmt.Println(m[2], m[4])
	fmt.Println(m[3])
	if m[3] == (CBC{Iv: "not empty string"}) {
		fmt.Println("Yes")
	}
	fmt.Println("No")
}

//func main() {
//	cbc := CBC{Iv: "6368616e676520746869732070617373"}
//	filename := "test1"
//	cbc.EncryptToFile(filename, "plaintexstMsg")
//	plaintext := cbc.DecryptFromFile(filename)
//	fmt.Println(plaintext)
//}

func (c *CBC) EncryptToFile(name, plaintext string) {
	ciphertextBytes := c.CBCEncrypter(plaintext)
	fmt.Println("cipherText (as bytes): ", ciphertextBytes)
	err := os.WriteFile(name, ciphertextBytes, 0644)
	check(err)
}

func (c *CBC) DecryptFromFile(name string) string {
	ciphertext, err := os.ReadFile(name)
	check(err)
	plaintextBytes := c.CBCDecrypter(ciphertext)
	plaintext := string(plaintextBytes[:])
	return plaintext
}
//https://golang.org/src/crypto/cipher/example_test.go
func (c *CBC) CBCDecrypter(ciphertext []byte) []byte {
	key, _ := hex.DecodeString(c.Iv)

	block, err := aes.NewCipher(key)
	check(err)

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)
	fmt.Println("Ciphertext 3495:", string(ciphertext[:]), ciphertext, "len(ciphertext)=", len(ciphertext))
	oldCiphertext := ciphertext[:]
	for i, _ := range oldCiphertext {
		reversedIndex := len(ciphertext)-1-i
		fmt.Println("tst", ciphertext[reversedIndex], string(ciphertext[reversedIndex]), "inex:", reversedIndex) // Suggestion: do `last := len(s)-1` before the loop
		if oldCiphertext[reversedIndex] == byte(0) {
			fmt.Println("removing index:", reversedIndex)
			fmt.Println(string(ciphertext[:]))
			ciphertext = RemoveIndex(ciphertext, reversedIndex)
		}
	}
	fmt.Println("asdf")
	fmt.Println("result:", string(ciphertext[:]))

	// If the original plaintext lengths are not a multiple of the block
	// size, padding would have to be added when encrypting, which would be
	// removed at this point. For an example, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. However, it's
	// critical to note that ciphertexts must be authenticated (i.e. by
	// using crypto/hmac) before being decrypted in order to avoid creating
	// a padding oracle.

	fmt.Printf("%s\n", ciphertext)
	// Output: exampleplaintext
	return ciphertext
}
func RemoveIndex(s []byte, index int) []byte {
	if len(s) == index+1 {
		return s[:index]
	}
	return append(s[:index], s[index+1:]...)
}

func (c *CBC) CBCEncrypter(plaintext string) []byte {
	plaintextBytes := []byte(plaintext)

	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. For an example of such padding, see https://tools.ietf.org/html/rfc5246#section-6.2.3.2
	remainder := len(plaintextBytes)%aes.BlockSize
	if remainder != 0 {
		// "plaintext is not a multiple of the block size"
		// We need to pad it so that it is
		bytesToPad := aes.BlockSize - remainder

		newPlaintextBytes := make([]byte, len(plaintextBytes)+bytesToPad)
		copy(newPlaintextBytes, plaintextBytes)
		copy(newPlaintextBytes[len(plaintextBytes):], bytes.Repeat([]byte{byte(0)}, bytesToPad) )
		plaintextBytes = newPlaintextBytes
	}

	key, _ := hex.DecodeString(c.Iv)
	block, err := aes.NewCipher(key)
	check(err)

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintextBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintextBytes)

	fmt.Printf("Encrypted ciphertext: '%x'\n", ciphertext)
	return ciphertext
}


func check(e error) {
	if e != nil {panic(e)}
}
