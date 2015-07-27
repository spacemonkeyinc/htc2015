package dh

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

const (
	prime = "00f2b2ab9d7b23c84f9f0ec2f3bc40c5c4ec4764a7c3d0144966" +
		"2620dd43f3d97a64515a2af5b3c8e3f224b8d18d07b6b6226120" +
		"0ad848f5ff8ac19a1b7343994de846de69c1c2ee5e62fe4ed374" +
		"e685e486f1b897d72d01df5c99ae72b8e9a31777ccaa11a5ae6c" +
		"a08cfc810269337660248d0be9b8214ecdd4656f207d2977a736" +
		"4e443acf431af76aead7224f86a03eb9998692acebd50c558ce9" +
		"a7fefc37ab242f0c19b51a0167d5dae94b853210f6f492a9bbb3" +
		"9ad809396b44a299bd85acafdfedbc4d21ae2ec307ab3dab09d7" +
		"99c6011c41cf813d621ef205cf2276d0cf7acf09108e14a8b8dd" +
		"e1ee2045deaebdb529dbd187d4ee4b30a94658b156ac33"
	Base = 2
)

var (
	Prime = mustParse(prime)
)

func mustParse(prime string) (rv *big.Int) {
	var success bool
	rv, success = new(big.Int).SetString(prime, 16)
	if !success {
		panic("bad prime")
	}
	return rv
}

func Private() (rv *big.Int) {
	var err error
	rv, err = rand.Int(rand.Reader, Prime)
	if err != nil {
		panic(err)
	}
	return rv
}

func Public(private *big.Int) *big.Int {
	return new(big.Int).Exp(big.NewInt(Base), private, Prime)
}

func SessionId(your_private, other_public *big.Int) [32]byte {
	sess := new(big.Int).Exp(other_public, your_private, Prime)
	width := len(prime) / 2
	b := sess.Bytes()
	if len(b) < width {
		b = append(make([]byte, width-len(b)), b...)
	}
	return sha256.Sum256(b)
}
