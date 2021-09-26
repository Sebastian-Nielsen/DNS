package main

import (
	"crypto/rand"
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

//func main() {
//	//fmt.Println("Result:", Encrypt(2, SecretKey{2, 2}))
//	n, d := KeyGen(20)
//	publicKey := PublicKey{N:n, E:big.NewInt(3)}
//	secretKey := SecretKey{N:n, D:d}
//	m := big.NewInt( 123456 )   // kan ikke klare beskeder med længde > 6 ?!?
//	fmt.Println("original:", m)
//	fmt.Println("encrypted:", Encrypt(m, secretKey))
//	fmt.Println("decrypted m:", Decrypt(Encrypt(m, secretKey), publicKey))
//	fmt.Println("original m:", m)
//}

func KeyGen(k int) (*big.Int, *big.Int) {
	n, p, q := compute_n(k)
	d := compute_d(p, q)

	if !( n.BitLen() == k ) {
		fmt.Println( "ERROR 8231")
	}
	return n, d
}
func Encrypt(m *big.Int, key SecretKey) *big.Int { // Compute the signature of m
	n := key.N
	d := key.D

	fmt.Println("m:", m)
	fmt.Println("d:", d)
	dRaisedToM := big.NewInt(0).Exp(m, d, nil)
	fmt.Println("d and m:", d, m)
	//fmt.Println("debug d^m ", dRaisedToM)
	c := big.NewInt(0).Mod(dRaisedToM, n)     // m^d % n
	return c
}
func Decrypt(c *big.Int, key PublicKey) *big.Int {
	n := key.N
	e := key.E

	eRaisedToC := big.NewInt(0).Exp(c, e, nil)
	m := big.NewInt(0).Mod(eRaisedToC, n)     // m^e % n
	return m
}
func compute_d(p, q *big.Int) *big.Int {

	qMin1 := big.NewInt(0).Sub(q, big.NewInt(1))
	pMin1 := big.NewInt(0).Sub(p, big.NewInt(1))

	//big.NewInt(int64(math.Pow(3, 1))),
	//https://youtu.be/Qgow8pVNjr0?t=561
	// d=e^(-1) mod (p-1)(q-1)    equivalent to   de=1 mod (p-1)(q-1)
	// e^(-1) is not as in the usual integer sense in the above
	d := big.NewInt(0).ModInverse(
		big.NewInt(3),
		big.NewInt(0).Mul(qMin1, pMin1),  // (p-1)(q-1)
	)
	fmt.Println()
	fmt.Println("q-1 =", qMin1)
	fmt.Println("p-1 =", pMin1)
	fmt.Println("d = 3^(-1) mod (p-1)(q-1) = d,   len(d)=",d.BitLen())
	fmt.Println()
	return d
}

func compute_n(k int) (*big.Int, *big.Int, *big.Int) {
	for {  // There is no do-while in golang, so just while forever until condition is meet then break
		p, err1 := rand.Prime(rand.Reader, k/2)
		if err1 != nil {fmt.Println("1", err1)}
		q, err2 := rand.Prime(rand.Reader, k/2)
		if err2 != nil {fmt.Println("2", err2)}
		n := big.NewInt(0).Mul(p, q)
		fmt.Println("p:", p, "q:", q)
		fmt.Println("n:", n)
		fmt.Println("length of p:", p.BitLen())
		fmt.Println("length of q:", q.BitLen())
		fmt.Println("length of n:", n.BitLen())
		qMin1 := big.NewInt(0).Sub(q, big.NewInt(1))
		pMin1 := big.NewInt(0).Sub(p, big.NewInt(1))
		gcd3AndPMin1 := big.NewInt(0).GCD(nil, nil, big.NewInt(3), pMin1)
		gcd3AndQMin1 := big.NewInt(0).GCD(nil, nil, big.NewInt(3), qMin1)
		gcd3AndPMin1Min1 :=  big.NewInt(0).Sub(gcd3AndPMin1, big.NewInt(1))
		if gcd3AndPMin1.Cmp(gcd3AndQMin1) == 0 &&
			gcd3AndPMin1Min1.Cmp(big.NewInt(0)) == 0  {
			fmt.Println("p and q does satisfy: gcd(3,p-1) = gcd(3,q-1) = 1\n-------------------------")
			return n, p, q
		} else {
			fmt.Println("p and q does not satisfy: gcd(3,p-1) = gcd(3,q-1) = 1\n")
		}
	}
}