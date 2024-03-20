package crypto

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"testing"
)

func TestSignerRSA_Sign(t *testing.T) {
	rsaGenerator := RSAGenerator{}
	keyPair, _ := rsaGenerator.Generate()

	dataToBeSigned := []byte("test message")
	signer := SignerRSA{KeyPair: keyPair}
	sign, err := signer.Sign(dataToBeSigned)
	if err != nil {
		t.Fail()
	}

	hashDataToBeSigned := sha256.New()
	hashDataToBeSigned.Write(dataToBeSigned)
	hashSumDataToBeSigned := hashDataToBeSigned.Sum(nil)

	err = rsa.VerifyPSS(keyPair.Public, crypto.SHA256, hashSumDataToBeSigned, sign, nil)
	if err != nil {
		t.Fail()
	}
}

func TestSignerECC_Sign(t *testing.T) {
	rsaGenerator := ECCGenerator{}
	keyPair, _ := rsaGenerator.Generate()

	dataToBeSigned := []byte("test message")
	signer := SignerECDSA{KeyPair: keyPair}
	_, err := signer.Sign(dataToBeSigned)
	if err != nil {
		t.Fail()
	}
}
