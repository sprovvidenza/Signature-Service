package persistence

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func createDevice(t *testing.T) domain.DeviceEntity {
	generator := crypto.RSAGenerator{}
	key, err := generator.Generate()
	if err != nil {
		t.Fail()
	}
	marshaler := crypto.RSAMarshaler{}
	pub, private, err := marshaler.Marshal(*key)
	if err != nil {
		t.Fail()
	}
	device := domain.DeviceEntity{Id: uuid.New(), Label: "test", PrivateKey: private, PubKey: pub, SignCounter: 0}
	return device
}

func TestDeviceInMemoryRepo_save(t *testing.T) {
	repo := DeviceInMemoryRepo{data: make(map[string]domain.DeviceEntity)}
	device := createDevice(t)

	saved := repo.Save(device)

	if !reflect.DeepEqual(saved, device) {
		t.Fail()
	}
}

func TestDeviceInMemoryRepo_findById(t *testing.T) {
	repo := DeviceInMemoryRepo{data: make(map[string]domain.DeviceEntity)}
	device := createDevice(t)
	repo.Save(device)
	findById, ok := repo.FindById(device.Id.String())
	if !ok {
		t.Fail()
	}
	if !reflect.DeepEqual(device, findById) {
		t.Fail()
	}
}

func TestDeviceInMemoryRepo_UpdateCounterById(t *testing.T) {
	repo := DeviceInMemoryRepo{data: make(map[string]domain.DeviceEntity)}
	device := createDevice(t)
	repo.Save(device)
	counter, _ := repo.UpdateCounterById(device.Id.String())
	if counter.SignCounter != 1 {
		t.Fail()
	}
	byId, _ := repo.FindById(device.Id.String())
	if byId.SignCounter != 1 {
		t.Fail()
	}

}
