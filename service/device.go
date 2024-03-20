package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
)

var mutex sync.Mutex

type DeviceService struct {
	signerService SignerService
	repo          persistence.DeviceRepo
}

func NewDeviceService(service SignerService, repo persistence.DeviceRepo) *DeviceService {

	return &DeviceService{signerService: service, repo: repo}
}

func (d *DeviceService) CreateSignatureDevice(algorithm string, label string) domain.CreateSignatureDeviceResponse {
	pub, private, err := d.signerService.createKeyPair(algorithm)
	device := domain.DeviceEntity{Id: uuid.New(), Label: label, PubKey: pub, PrivateKey: private, SignCounter: 0, Algorithm: algorithm}
	saved := d.repo.Save(device)
	if err != nil {
		return domain.CreateSignatureDeviceResponse{}
	}
	return domain.CreateSignatureDeviceResponse{Id: saved.Id.String(), Label: saved.Label, SignatureCount: saved.SignCounter}
}

func (d *DeviceService) GetAllDevice() []domain.AllDevice {
	deviceEntities := d.repo.FindAll()
	var values []domain.AllDevice
	for _, element := range deviceEntities {
		device := domain.AllDevice{Id: element.Id.String(), Label: element.Label, Counter: element.SignCounter}
		values = append(values, device)
	}
	return values
}

func (d *DeviceService) SignTransaction(deviceId string, data string) (domain.SignatureResponse, error) {
	device, ok := d.repo.FindById(deviceId)
	if !ok {
		return domain.SignatureResponse{}, errors.New("NOT_FOUND")
	}
	signer := d.signerService.getSigner(device)
	if signer == nil {
		return domain.SignatureResponse{}, errors.New("NOT_ALLOWED")
	}
	sign, _ := signer.Sign([]byte(data))
	mutex.Lock()
	log.Printf("Start transaction \n")
	signature := base64.StdEncoding.EncodeToString(sign)
	signedData := d.makeSignedData(device, data)
	device, _ = d.repo.UpdateCounterById(deviceId)
	mutex.Unlock()
	log.Printf("Stop transaction \n")
	return domain.SignatureResponse{Signature: signature, SignedData: signedData}, nil
}

func (d *DeviceService) makeSignedData(device domain.DeviceEntity, data string) string {
	var lastSignB64 string
	if device.SignCounter == 0 {
		lastSignB64 = base64.StdEncoding.EncodeToString([]byte(device.Id.String()))
	} else {
		lastSignB64 = time.Now().String()
	}
	signedData := fmt.Sprintf("%x", device.SignCounter) + "_" + data + "_" + lastSignB64
	return signedData
}
