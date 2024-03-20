package api

import (
	"encoding/json"
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
	"io"
	"net/http"
)

// TODO: REST endpoints ...

// Create device with PUT and retrieve all device with GET
func (s *Server) CreateDevice(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPut:
		createDevice(response, request, s.deviceService)
	case http.MethodGet:
		getAllDevice(response, request, s.deviceService)
	default:
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			request.Method + " Not Allowed",
		})
		return
	}
}
func getAllDevice(response http.ResponseWriter, request *http.Request, deviceService service.DeviceService) {
	devices := deviceService.GetAllDevice()
	WriteAPIResponse(response, http.StatusOK, devices)
}

func createDevice(response http.ResponseWriter, request *http.Request, s service.DeviceService) {
	var body domain.CreateDeviceBody
	bodyByte, _ := io.ReadAll(request.Body)
	err := json.Unmarshal(bodyByte, &body)
	if err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"input body il malformed",
		})
		return
	}

	if body.Algorithm == "" {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"alg is mandatory",
		})
		return
	}

	signatureResponse := s.CreateSignatureDevice(body.Algorithm, body.Label)

	WriteAPIResponse(response, http.StatusCreated, signatureResponse)
}

// Sign device
func (s *Server) SignDevice(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			request.Method + " Not Allowed",
		})
		return
	}
	var body domain.SignDeviceBody
	bodyByte, _ := io.ReadAll(request.Body)
	err := json.Unmarshal(bodyByte, &body)
	if err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"input body il malformed",
		})
		return
	}

	if body.DeviceId == "" {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"device_id is mandatory",
		})
		return
	}
	if body.DataToSign == "" {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"data_to_sign is mandatory",
		})
		return
	}

	signatureResponse, err := s.deviceService.SignTransaction(body.DeviceId, body.DataToSign)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			WriteErrorResponse(response, http.StatusNotFound, []string{
				fmt.Sprintf("device with id %s not found ", body.DeviceId),
			})
			return
		case "NOT_ALLOWED":
			WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
				fmt.Sprintf("The signature in device with id %s is not allowed ", body.DeviceId),
			})
			return
		}

	}

	WriteAPIResponse(response, http.StatusCreated, signatureResponse)
}
