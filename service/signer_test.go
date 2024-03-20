package service

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"reflect"
	"testing"
)

func TestSignerService_createKeyPair(t *testing.T) {
	methodMap := make(map[string]func() ([]byte, []byte, error))
	methodMap["RSA"] = RsaGenerator
	methodMap["ECDSA"] = EccGenerator

	service := SignerService{signerGenerators: methodMap}

	_, _, err := service.createKeyPair("RSA")
	if err != nil {
		t.Errorf("%e", err)
	}

	_, _, err = service.createKeyPair("ECDSA")
	if err != nil {
		t.Errorf("%e", err)
	}

	_, _, err = service.createKeyPair("NOT_EXIST")
	if err == nil {
		t.Errorf("Should go into error")
	}

}

func TestSignerService_getSigner(t *testing.T) {
	signerMap := make(map[string]func(entity domain.DeviceEntity) crypto.Signer)
	signerMap["RSA"] = RsaSigner
	signerMap["ECDSA"] = EccSigner

	service := SignerService{signer: signerMap}
	rsaGen := crypto.RSAGenerator{}
	rsaKeypair, _ := rsaGen.Generate()
	rsaMarshaler := crypto.RSAMarshaler{}
	rsaPub, rsaPriv, _ := rsaMarshaler.Marshal(*rsaKeypair)
	signerRsa := service.getSigner(domain.DeviceEntity{Algorithm: "RSA", PrivateKey: rsaPriv, PubKey: rsaPub})
	if signerRsa == nil {
		t.Fail()
	}
	eccGen := crypto.ECCGenerator{}
	eccKeypair, _ := eccGen.Generate()
	eccMarshaler := crypto.ECCMarshaler{}
	eccPub, eccPriv, _ := eccMarshaler.Encode(*eccKeypair)
	signerEcc := service.getSigner(domain.DeviceEntity{Algorithm: "ECDSA", PubKey: eccPub, PrivateKey: eccPriv})
	if signerEcc == nil {
		t.Fail()
	}
	typeOfRsa := reflect.TypeOf(signerRsa).String()

	typeOfEcc := reflect.TypeOf(signerEcc).String()

	if typeOfRsa != "*crypto.SignerRSA" {
		t.Fail()
	}
	if typeOfEcc != "*crypto.SignerECDSA" {
		t.Fail()
	}
}
