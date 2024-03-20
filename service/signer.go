package service

import (
	"errors"
	"fmt"
	"github.com/sprovvidenza/Signature-Service/crypto"
	"github.com/sprovvidenza/Signature-Service/domain"
	"log"
)

type SignerService struct {
	signerGenerators map[string]func() ([]byte, []byte, error)
	signer           map[string]func(entity domain.DeviceEntity) crypto.Signer
}

func NewSignerService(signerGenerators map[string]func() ([]byte, []byte, error), signer map[string]func(entity domain.DeviceEntity) crypto.Signer) SignerService {

	return SignerService{signer: signer, signerGenerators: signerGenerators}
}

// generic method to create a key pair of a particular algorithm, Only supported
func (s *SignerService) createKeyPair(algorithm string) ([]byte, []byte, error) {
	f, ok := s.signerGenerators[algorithm]
	if ok {
		return f()
	}
	return nil, nil, errors.New(fmt.Sprintf("Method %s is not supported", algorithm))
}

// generic method to provide a signer, Only supported
func (s *SignerService) getSigner(entity domain.DeviceEntity) crypto.Signer {
	f, ok := s.signer[entity.Algorithm]
	if ok {
		return f(entity)
	}
	return nil
}

func RsaSigner(entity domain.DeviceEntity) crypto.Signer {
	rsaMarshaler := crypto.RSAMarshaler{}
	rsaKeyPair, _ := rsaMarshaler.Unmarshal(entity.PrivateKey)
	rsa := crypto.SignerRSA{KeyPair: rsaKeyPair}
	log.Printf("Using RSA signer")
	return &rsa
}

func EccSigner(entity domain.DeviceEntity) crypto.Signer {
	eccMarshaler := crypto.ECCMarshaler{}
	eccKeyPair, _ := eccMarshaler.Decode(entity.PrivateKey)
	rsa := crypto.SignerECDSA{KeyPair: eccKeyPair}
	log.Printf("Using ECDSA signer")
	return &rsa
}

func RsaGenerator() ([]byte, []byte, error) {
	generator := crypto.RSAGenerator{}
	generate, err := generator.Generate()
	if err != nil {
		return nil, nil, err
	}
	marshaler := crypto.RSAMarshaler{}
	public, private, err := marshaler.Marshal(*generate)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("Generate RSA key pair")
	return public, private, nil
}

func EccGenerator() ([]byte, []byte, error) {
	generator := crypto.ECCGenerator{}
	generate, err := generator.Generate()
	if err != nil {
		return nil, nil, err
	}
	marshaler := crypto.ECCMarshaler{}
	public, private, err := marshaler.Encode(*generate)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("Generate ECDSA key pair")
	return public, private, nil
}
