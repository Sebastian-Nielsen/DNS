package Cryptography

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type PublicKey struct {
	N *big.Int
	E *big.Int
}
type SecretKey struct {
	N *big.Int
	D *big.Int
}

func (p *PublicKey) ToString() string {
	return p.N.String() + ":" + p.E.String()
}


//func Test2() {
//	//fmt.Println("Result:", Encrypt(2, SecretKey{2, 2}))
//	n, d := KeyGen(18)
//	publicKey := PublicKey{N:n, E:big.NewInt(3)}
//	secretKey := SecretKey{N:n, D:d}
//	m := big.NewInt( 111116 )   // kan ikke klare beskeder med længde > 6 ?!?
//	fmt.Println("original:", m)
//	fmt.Println("encrypted:", Encrypt(m, publicKey))
//	decryptedM := Decrypt(Encrypt(m, publicKey), secretKey)
//	//decrypted_m := Encrypt(Decrypt(m, publicKey), secretKey)
//	fmt.Println("---------------------")
//	fmt.Println("Success?   ", decryptedM.Cmp(m) == 0)
//	fmt.Println("decrypted m:", decryptedM)
//	fmt.Println("original m:", m)
//	fmt.Println("---------------------")
//}


func GetKeys(k int) (PublicKey, SecretKey) {
	n, d := keyGen(k)
	pk := PublicKey{N:n, E:big.NewInt(3)}
	sk := SecretKey{N:n, D:d}
	return pk, sk
}

func keyGen(k int) (*big.Int, *big.Int) {
	n, p, q := compute_n(k)
	d := compute_d(p, q)

	if !( n.BitLen() == k ) {
		fmt.Println( "ERROR 8231")
	}
	return n, d
}
func Encrypt(m *big.Int, key PublicKey) *big.Int { // Compute the signature of m
	n := key.N
	e := key.E

	//fmt.Println("m:", m)
	//fmt.Println("e:", e)
	c := big.NewInt(0).Exp(m, e, n)
	// fmt.Println("Debug m^e % n =", c)
	return c
}
func Decrypt(c *big.Int, key SecretKey) *big.Int {
	n := key.N
	d := key.D

	m := big.NewInt(0).Exp(c, d, n)
	// fmt.Println("Debug2 c^d % n =", m)
	return m
}
func compute_d(p, q *big.Int) *big.Int {
	qMin1 := big.NewInt(0).Sub(q, big.NewInt(1))
	pMin1 := big.NewInt(0).Sub(p, big.NewInt(1))

	//qMin1 := big.NewInt(0).Sub(q, big.NewInt(1))
	//pMin1 := big.NewInt(0).Sub(p, big.NewInt(1))
	//big.NewInt(int64(math.Pow(3, 1))),
	//https://youtu.be/Qgow8pVNjr0?t=561
	// d=e^(-1) mod (p-1)(q-1)    equivalent to   de=1 mod (p-1)(q-1)
	// e^(-1) is not as in the usual integer sense in the above          3^(-1)  mod (q-1)(p-1)
	d := big.NewInt(0).ModInverse(
		big.NewInt(3),
		big.NewInt(0).Mul(qMin1, pMin1),  // (p-1)(q-1)
	)
	//fmt.Println()
	//fmt.Println("q-1 =", qMin1)
	//fmt.Println("p-1 =", pMin1)
	//fmt.Println("(p-1)(q-1) =", big.NewInt(0).Mul(pMin1, qMin1))
	//fmt.Println("d = 3^(-1) mod (p-1)(q-1) = d,   len(d)=",d.BitLen())
	//fmt.Println("d = 3^(-1) mod (p-1)(q-1) =", d)
	//fmt.Println()
	return d
}

func compute_n(k int) (*big.Int, *big.Int, *big.Int) {
	for {  // There is no do-while in golang, so just while forever until condition is meet then break
		p, err1 := rand.Prime(rand.Reader, k/2)
		if err1 != nil {fmt.Println("1", err1)}
		q, err2 := rand.Prime(rand.Reader, k/2)
		if err2 != nil {fmt.Println("2", err2)}
		n := big.NewInt(0).Mul(p, q)
		//fmt.Println("p:", p, "q:", q)
		//fmt.Println("n:", n)
		//fmt.Println("length of p:", p.BitLen())
		//fmt.Println("length of q:", q.BitLen())
		//fmt.Println("length of n:", n.BitLen())
		qMin1 := big.NewInt(0).Sub(q, big.NewInt(1))
		pMin1 := big.NewInt(0).Sub(p, big.NewInt(1))
		gcd3AndPMin1 := big.NewInt(0).GCD(nil, nil, big.NewInt(3), pMin1)
		gcd3AndQMin1 := big.NewInt(0).GCD(nil, nil, big.NewInt(3), qMin1)
		if gcd3AndPMin1.Cmp(gcd3AndQMin1) == 0 &&
			gcd3AndPMin1.Cmp(big.NewInt(1)) == 0 &&
			p.Cmp(q) != 0 {
			//fmt.Println("p and q does satisfy: gcd(3,p-1) = gcd(3,q-1) = 1\n-------------------------")
			return n, p, q
		// } else {
			// fmt.Println("p and q does not satisfy: gcd(3,p-1) = gcd(3,q-1) = 1\n ")
		}
	}
}

func BigInt_verify(signature *big.Int, msg *big.Int, pk PublicKey) bool {
	hashedMsg := new(big.Int)
	hashedMsg.SetBytes(Hash(msg))
	unsignedMsg := Encrypt(signature, pk)
	//fmt.Println("Verify2 ----")
	//fmt.Println("msg: " + msg.String() + "\nhashed: " + hashedMsg.String())
	//fmt.Println("bigint: " + unsignedMsg.String() + "\nsignature: " + signature.String())
	return hashedMsg.Cmp(unsignedMsg) == 0
}

func Verify(signature string, msg string, pk PublicKey) bool {
	// convert {signature} of type string to big.Int
	/*sigByteArr := []byte(signature)
	sigBigInt := new(big.Int)
	sigBigInt.SetBytes(sigByteArr)*/

	sigBigInt := new(big.Int)
	sigBigInt.SetString(signature, 10)


	// convert {msg} of type string to big.Int
	msgByteArr := []byte(msg)
	msgBigInt := new(big.Int)
	msgBigInt.SetBytes(msgByteArr)


	//fmt.Println("Vefify1 ----")
	//fmt.Println("msg: " + msg + "\nbigint: " + msgBigInt.String())
	//fmt.Println("signature: " + signature + "\nbigint: " + sigBigInt.String())
	return BigInt_verify(sigBigInt, msgBigInt, pk)
}

func Hash(msg *big.Int) []byte {
	sha := sha256.New()
	sha.Write([]byte(msg.String()))
	return sha.Sum(nil)
}

func BigInt_createSignature(msg *big.Int, sk SecretKey) *big.Int {
	hashedMsg := new(big.Int)
	hashedMsg.SetBytes(Hash(msg))
	//fmt.Println("sign2 ----")
	//fmt.Println("msg: " + msg.String() + "\nhashed: " + hashedMsg.String())
	return Decrypt(hashedMsg, sk)
}

func CreateSignature(msg string, sk SecretKey) string {
	byteArr := []byte(msg)
	bigInt := new(big.Int)
	bigInt.SetBytes(byteArr)


	//fmt.Println("Sign ----")
	//fmt.Println("msg: " + msg + "\nbigint: " + bigInt.String())
	//fmt.Println("after: " + BigInt_createSignature(bigInt, sk).String())
	return BigInt_createSignature(bigInt, sk).String()
}
