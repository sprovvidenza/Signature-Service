package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// TODO: implement RSA and ECDSA signing ...

type SignerECDSA struct {
	KeyPair *ECCKeyPair
}

func (sEcc *SignerECDSA) Sign(dataToBeSigned []byte) ([]byte, error) {
	hashDataToBeSigned := sha256.New()
	_, err := hashDataToBeSigned.Write(dataToBeSigned)
	if err != nil {
		return nil, err
	}

	hashSumDataToBeSigned := hashDataToBeSigned.Sum(nil)
	r, s, err := ecdsa.Sign(rand.Reader, sEcc.KeyPair.Private, hashSumDataToBeSigned)
	if err != nil {
		return nil, err
	}
	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)

	return signature, nil
}

type SignerRSA struct {
	KeyPair *RSAKeyPair
}

func (sRsa *SignerRSA) Sign(dataToBeSigned []byte) ([]byte, error) {
	hashDataToBeSigned := sha256.New()
	_, err := hashDataToBeSigned.Write(dataToBeSigned)
	if err != nil {
		return nil, err
	}

	hashSumDataToBeSigned := hashDataToBeSigned.Sum(nil)
	signature, err := rsa.SignPSS(rand.Reader, sRsa.KeyPair.Private, crypto.SHA256, hashSumDataToBeSigned, nil)
	if err != nil {
		return nil, err
	}

	return signature, nil
}
