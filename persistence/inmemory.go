package persistence

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"log"
)

type DeviceRepo interface {
	Save(device domain.DeviceEntity) domain.DeviceEntity
	FindById(id string) (domain.DeviceEntity, bool)
	UpdateCounterById(id string) (domain.DeviceEntity, bool)
	FindAll() []domain.DeviceEntity
}

// TODO: in-memory persistence ...
type DeviceInMemoryRepo struct {
	data map[string]domain.DeviceEntity
}

func NewDeviceInMemoryRepo() DeviceInMemoryRepo {
	d := make(map[string]domain.DeviceEntity)
	return DeviceInMemoryRepo{data: d}
}

func (rep *DeviceInMemoryRepo) Save(device domain.DeviceEntity) domain.DeviceEntity {
	rep.data[device.Id.String()] = device
	log.Printf("Save device %s \n", device.Id)
	return rep.data[device.Id.String()]
}

func (rep *DeviceInMemoryRepo) FindById(id string) (domain.DeviceEntity, bool) {
	device, ok := rep.data[id]
	log.Printf("Get device %s \n", device.Id)
	return device, ok
}

func (rep *DeviceInMemoryRepo) FindAll() []domain.DeviceEntity {
	var values []domain.DeviceEntity
	for _, v := range rep.data {
		values = append(values, v)
	}
	log.Printf("Retrive all device")
	return values
}

func (rep *DeviceInMemoryRepo) UpdateCounterById(id string) (domain.DeviceEntity, bool) {
	device, ok := rep.data[id]
	device.SignCounter += 1
	rep.data[id] = device
	log.Printf("Update counter for device %s at %s \n", id, device.SignCounter)
	return device, ok
}
