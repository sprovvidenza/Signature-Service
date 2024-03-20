package api

import (
	"github.com/sprovvidenza/Signature-Service/crypto"
	"github.com/sprovvidenza/Signature-Service/domain"
	"github.com/sprovvidenza/Signature-Service/service"
)

// Generic interface to configure the algorithm supported
type SignConfiguration interface {
	Config() (map[string]func(entity domain.DeviceEntity) crypto.Signer, map[string]func() ([]byte, []byte, error))
}

// Implementation of config in memory
type InMemorySignConfiguration struct {
}

func (c *InMemorySignConfiguration) Config() (map[string]func(entity domain.DeviceEntity) crypto.Signer, map[string]func() ([]byte, []byte, error)) {
	signerMap := make(map[string]func(entity domain.DeviceEntity) crypto.Signer)
	signerMap["RSA"] = service.RsaSigner
	signerMap["ECDSA"] = service.EccSigner

	generatorsMap := make(map[string]func() ([]byte, []byte, error))
	generatorsMap["RSA"] = service.RsaGenerator
	generatorsMap["ECDSA"] = service.EccGenerator
	return signerMap, generatorsMap
}
