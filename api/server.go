package api

import (
	"encoding/json"
	"github.com/sprovvidenza/Signature-Service/persistence"
	"github.com/sprovvidenza/Signature-Service/service"
	"net/http"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress string
	deviceService service.DeviceService
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string) *Server {
	configuration := InMemorySignConfiguration{}
	signer, generators := configuration.Config()

	signerService := service.NewSignerService(generators, signer)
	repo := persistence.NewDeviceInMemoryRepo()
	deviceService := service.NewDeviceService(signerService, &repo)

	return &Server{
		listenAddress: listenAddress,
		// TODO: add services / further dependencies here ...
		deviceService: *deviceService,
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	mux := http.NewServeMux()

	mux.Handle("/api/v0/health", http.HandlerFunc(s.Health))
	mux.Handle("/api/v0/device", http.HandlerFunc(s.CreateDevice))
	mux.Handle("/api/v0/device/sign", http.HandlerFunc(s.SignDevice))

	// TODO: register further HandlerFuncs here ...

	return http.ListenAndServe(s.listenAddress, mux)
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}
