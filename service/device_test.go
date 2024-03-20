package service

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"testing"
)

func TestDeviceService_CreateSignatureDevice(t *testing.T) {
	methodMap := make(map[string]func() ([]byte, []byte, error))
	methodMap["RSA"] = RsaGenerator
	methodMap["ECDSA"] = EccGenerator
	repo := persistence.New()
	deviceService := DeviceService{signerService: SignerService{signerGenerators: methodMap}, repo: repo}
	signatureDevice := deviceService.CreateSignatureDevice("RSA", "Test")
	if signatureDevice == (CreateSignatureDeviceResponse{}) {
		t.Fail()
	}
}
