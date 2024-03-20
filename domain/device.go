package domain

import "github.com/google/uuid"

// TODO: signature device domain model ...

type DeviceEntity struct {
	Id          uuid.UUID
	Label       string
	PrivateKey  []byte
	PubKey      []byte
	SignCounter int
	Algorithm   string
}

type CreateSignatureDeviceResponse struct {
	Id             string `json:"id"`
	Label          string `json:"label"`
	SignatureCount int    `json:"signatureCount"`
}

type AllDevice struct {
	Id      string `json:"id"`
	Label   string `json:"label"`
	Counter int    `json:"counter"`
}

type SignatureResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

type CreateDeviceBody struct {
	Algorithm string `json:"alg"`
	Label     string `json:"label,omitempty"`
}

type SignDeviceBody struct {
	DataToSign string `json:"data_to_sign"`
	DeviceId   string `json:"device_id"`
}
