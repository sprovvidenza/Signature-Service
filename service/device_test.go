package service

import (
	"github.com/sprovvidenza/Signature-Service/domain"
	"github.com/sprovvidenza/Signature-Service/persistence"
	"testing"
)

func TestDeviceService_CreateSignatureDevice(t *testing.T) {
	methodMap := make(map[string]func() ([]byte, []byte, error))
	methodMap["RSA"] = RsaGenerator
	methodMap["ECDSA"] = EccGenerator
	repo := persistence.NewDeviceInMemoryRepo()
	deviceService := DeviceService{signerService: SignerService{signerGenerators: methodMap}, repo: &repo}
	signatureDevice := deviceService.CreateSignatureDevice("RSA", "Test")
	if signatureDevice == (domain.CreateSignatureDeviceResponse{}) {
		t.Fail()
	}
}
