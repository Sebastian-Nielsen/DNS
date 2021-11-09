package Cryptography

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

/*type Signature struct {  asdfasdf
	Value string
}
func (s *Signature) ToString() string {
	return s.Value
}*/


/*	          Source: https://stackoverflow.com/a/5142164/7123519
	^                         Start anchor
	(?=.*[A-Z].*[A-Z])        Ensure string has two uppercase letters.
	(?=.*[!@#$&*])            Ensure string has one special case letter.
	(?=.*[0-9].*[0-9])        Ensure string has two digits.
	(?=.*[a-z].*[a-z].*[a-z]) Ensure string has three lowercase letters.
	.{8}                      Ensure string is of length 8.
	$                         End anchor.
*/
func isPasswordStrongEnough(password string) bool {
	tests := []string{".*[A-Z].*[A-Z]", ".*[!@#$&*]", ".*[0-9].*[0-9]", ".*[a-z].*[a-z].*[a-z]", ".{8,}"}
	for _, test := range tests {
		isSecure, _ := regexp.MatchString(test, password)
		if !isSecure {
			return false
		}
	}
	return true
}
func Generate(filename string, password string) (PublicKey, error) {
	if !isPasswordStrongEnough(password) {
		return PublicKey{}, errors.New("password too weak! (atleast two uppercase letters, one special letter, two digits, three lowercase letters and a length of 8)")
	}
	pk, sk := GenKeys(2000)

	// We assume that the hashed password byte array has a length of exactly 60
	hashedPw, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	ctrEncrypter := CTR{SecretKey: string(hashedPw[:32])}
	ctrEncrypter.EncryptToFile(filename, string(hashedPw) + sk.ToString())

	return pk, nil
}
func Sign(filename string, password string, msg []byte) (string, error) {
	hashedPw, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	ctrEncrypter := CTR{SecretKey: string(hashedPw[:32])}
	arrOfSks := ctrEncrypter.DecryptFromFile(filename)
	RSAsk := ToSecretKey(arrOfSks[61:])


	isValidPassword := bcrypt.CompareHashAndPassword(hashedPw, []byte(password)) == nil
	if isValidPassword {
		return CreateSignature(string(msg), RSAsk), nil
	}
	return "", errors.New("wrong password for the file")
}

func main() {
	fmt.Println("test")
}
