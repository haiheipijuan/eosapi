package ecc

import (
	cryptorand "crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
)

func NewRandomPrivateKey() (*PrivateKey, error) {
	return newRandomPrivateKey(cryptorand.Reader)
}

func NewDeterministicPrivateKey(randSource io.Reader) (*PrivateKey, error) {
	return newRandomPrivateKey(randSource)
}

func newRandomPrivateKey(randSource io.Reader) (*PrivateKey, error) {
	rawPrivKey := make([]byte, 32)
	written, err := io.ReadFull(randSource, rawPrivKey)
	if err != nil {
		return nil, fmt.Errorf("error feeding crypto-rand numbers to seed ephemeral private key: %s", err)
	}
	if written != 32 {
		return nil, fmt.Errorf("couldn't write 32 bytes of randomness to seed ephemeral private key")
	}

	h := sha256.New()
	h.Write(rawPrivKey)
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), h.Sum(nil))

	return &PrivateKey{privKey: privKey}, nil
}

func NewPrivateKey(wif string) (*PrivateKey, error) {
	wifObj, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{privKey: wifObj.PrivKey}, nil
}

type PrivateKey struct {
	privKey *btcec.PrivateKey
}

func (p *PrivateKey) PublicKey() PublicKey {
	return PublicKey(p.privKey.PubKey().SerializeCompressed())
}

// Sign signs a 32 bytes SHA256 hash..
func (p *PrivateKey) Sign(hash []byte) (Signature, error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash should be 32 bytes")
	}

	sig, err := p.privKey.Sign(hash)

	// der := sig.Serialize()
	// lenR := der[3]
	// lenS := der[5+lenR]
	// i := 0
	// if lenR == 32 && lenS == 32 {
	// }

	buf := make([]byte, 65, 65)
	//buf[0] =
	copy(buf[1:33], sig.R.Bytes())
	copy(buf[33:], sig.S.Bytes())

	checksum := ripemd160checksum(buf)
	textSig := base58.Encode(append(buf, checksum[:4]...))
	fmt.Println("OUTPUT sig", "EOS"+textSig)

	compactSig, err := btcec.SignCompact(btcec.S256(), p.privKey, hash, true)
	if err != nil {
		return nil, err
	}

	return Signature(compactSig), nil
}

func (p *PrivateKey) String() string {
	wif, _ := btcutil.NewWIF(p.privKey, &chaincfg.Params{PrivateKeyID: '\x80'}, false) // no error possible
	return wif.String()
}
