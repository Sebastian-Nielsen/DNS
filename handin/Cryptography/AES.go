package cryptography

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


// func main() {
// 	cbc := CBC{Iv: "6368616e676520746869732070617373"}
// 	filename := "test1"
// 	msg := "plaintexstMsg"
// 	fmt.Printf("msg: %s\n", msg)
// 	cbc.EncryptToFile(filename, msg)
// 	plaintext := cbc.DecryptFromFile(filename)
// 	fmt.Printf("en- and decrypted: %s\nEqual: %t", plaintext, msg==plaintext)
// }

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

	// Remove padding
	// fmt.Println(ciphertext)
	paddedBytes := int(ciphertext[len(ciphertext)-1])
	// fmt.Println(paddedBytes)
	ciphertext = ciphertext[:len(ciphertext)-paddedBytes]
	// fmt.Println(ciphertext)

	// fmt.Println("Ciphertext 3495:", ciphertext[:])
	// fmt.Println("Ciphertext 3495:", string(ciphertext))
	// fmt.Println("Ciphertext 3495:", []byte(string(ciphertext)))


	// If the original plaintext lengths are not a multiple of the block
	// size, padding would have to be added when encrypting, which would be
	// removed at this point. For an example, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. However, it's
	// critical to note that ciphertexts must be authenticated (i.e. by
	// using crypto/hmac) before being decrypted in order to avoid creating
	// a padding oracle.

	// fmt.Printf("%s\n", ciphertext)
	// Output: exampleplaintext
	return ciphertext
}

func (c *CBC) CBCEncrypter(plaintext string) []byte {
	plaintextBytes := []byte(plaintext)
	// fmt.Println("Bytes:")
	// fmt.Println(plaintextBytes)
	// fmt.Println("")



	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. For an example of such padding, see https://tools.ietf.org/html/rfc5246#section-6.2.3.2
	remainder := len(plaintextBytes)%aes.BlockSize
	// if remainder == 0 { remainder = aes.BlockSize }
		// "plaintext is not a multiple of the block size"
		// We need to pad it so that it is
	bytesToPad := aes.BlockSize - remainder

		// newPlaintextBytes := make([]byte, bytesToPad)
		// fmt.Println(append(newPlaintextBytes, plaintextBytes...))
		// copy(plaintextBytes, append(newPlaintextBytes, plaintextBytes...))
		// plaintextBytes = append(newPlaintextBytes, plaintextBytes...)
		
		
		// newPlaintextBytes := make([]byte, len(plaintextBytes)+bytesToPad)
		// copy(newPlaintextBytes, plaintextBytes)
		// fmt.Println(plaintextBytes)
		// fmt.Println(newPlaintextBytes)
		// copy(newPlaintextBytes[len(plaintextBytes):], bytes.Repeat([]byte{byte(bytesToPad)}, bytesToPad) )
		// plaintextBytes = newPlaintextBytes
		
	paddingBytes := bytes.Repeat([]byte{byte(bytesToPad)}, bytesToPad)
	plaintextBytes = append(plaintextBytes, paddingBytes...)
	// fmt.Println(plaintextBytes)
	
	// fmt.Println("Bytes padded:")
	// fmt.Println(plaintextBytes)
	// fmt.Println("")

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
