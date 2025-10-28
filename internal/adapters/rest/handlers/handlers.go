package handlers

import (
	"encoding/json"
	"fmt"
	"mockLogin/internal/config"
	"net/http"
	"time"

	domainPuesto "mockLogin/internal/domain/puesto"
	domainRequest "mockLogin/internal/domain/request"

	"github.com/golang-jwt/jwt/v4"
)

const (
	AVAILABLE_UPDATE = false

	RESPONSE_OK    = "ok"
	RESPONSE_ERROR = "failed"

	HS_TOKEN = 48 // revisar
)

var indicesUsados []int
var JWT = ""
var START time.Time
var indiceMapa int = 1

type Handlers struct {
	config   *config.Config                         // configuraciones
	Handlers map[string]map[string]http.HandlerFunc // handlers cargados
}

// NewHandlers inicializa los handlers con los servicios necesarios
func NewHandlers(cfg *config.Config) *Handlers {
	h := &Handlers{
		config:   cfg,
		Handlers: make(map[string]map[string]http.HandlerFunc),
	}

	// registra handlers necesarios - los que se necesiten
	h.RegisterHandler("GET", "/token", h.GetToken())
	h.RegisterHandler("POST", "/token/refresh", h.RefreshToken())
	h.RegisterHandler("POST", "/token/revoke", h.Bye())

	return h
}

// registrar dinámicamente handlers con métodos
func (h *Handlers) RegisterHandler(method string, endopoint string, handler http.HandlerFunc) {
	if _, exists := h.Handlers[method]; !exists {
		h.Handlers[method] = make(map[string]http.HandlerFunc)
	}

	// agregado de handler por método y ruta
	h.Handlers[method][endopoint] = handler
}

// /getToken
func (h *Handlers) GetToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if JWT == "" {
			// indice de mapa, indica en qué puesto está operando
			if indiceMapa > len(h.config.Puestos) {
				dataErr := domainRequest.MockLoginErrorData{
					Error: "No hay más puestos disponibles.",
				}
				responseJson(w, RESPONSE_ERROR, dataErr)
				return
			}
			currentPuesto := h.config.Puestos[indiceMapa]
			fmt.Printf("\nEste es el puesto actual: %+v\n", currentPuesto)

			var err error
			// genero JWT con valores del puesto
			JWT, START, err = generateJWT(h.config.Puestos[indiceMapa], 24*time.Hour, h.config.Server.SecretKey)
			if err != nil {
				data := domainRequest.MockLoginErrorData{
					Error: err.Error(),
				}

				responseJson(w, RESPONSE_ERROR, data)
			}

			data := domainRequest.MockLoginData{
				Jwt:   JWT,
				Start: START,
			}

			fmt.Printf("Token obtenido: %s\n", JWT)

			responseJson(w, RESPONSE_OK, data)
		} else {
			data := domainRequest.MockLoginData{
				Jwt:   JWT,
				Start: START,
			}

			fmt.Printf("Token obtenido: %s\n", JWT)

			responseJson(w, RESPONSE_OK, data)
		}
	}
}

// /refreshToken
func (h *Handlers) RefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if JWT == "" {
			dataErr := domainRequest.MockLoginErrorData{
				Error: "JWT not found",
			}
			responseJson(w, RESPONSE_ERROR, dataErr)
			return
		}
		fmt.Printf("\nToken actual: %s\n", JWT)

		currentPuesto := h.config.Puestos[indiceMapa]
		fmt.Printf("Refrescando token para el puesto: %+v\n", currentPuesto)

		var err error
		// genera nuevamenete el token para reasignar
		JWT, START, err = generateJWT(h.config.Puestos[indiceMapa], 24*time.Hour, h.config.Server.SecretKey)
		if err != nil {
			data := domainRequest.MockLoginErrorData{
				Error: err.Error(),
			}

			responseJson(w, RESPONSE_ERROR, data)
		}

		data := domainRequest.MockLoginData{
			Jwt:   JWT,
			Start: START,
		}

		fmt.Printf("Token Nuevo: %s\n", JWT)

		responseJson(w, RESPONSE_OK, data)
	}
}

// /bye
func (h *Handlers) Bye() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indiceMapa++

		// incrementa indice, así al consultar nuevamente opera otro puesto
		if indiceMapa > len(h.config.Puestos) {
			indiceMapa = 1
		}

		JWT = ""

		fmt.Printf("\nPuesto skipeado\n")

		responseJson(w, RESPONSE_OK, "Bye")
	}
}

// JWT
func generateJWT(puesto domainPuesto.Puesto, expiration time.Duration, secretKey string) (string, time.Time, error) {
	expirationTime := time.Now().Add(expiration)

	fmt.Printf("Esto es puesto: %+v\n", puesto)

	claims := jwt.MapClaims{
		"agencyBranch": puesto.AgencyBranch,
		"agencyCode":   puesto.AgencyCode,
		"workstation":  puesto.Workstation,
		"exp":          expirationTime.Unix(), // Tiempo de expiración en formato UNIX
		"iat":          time.Now().Unix(),     // Tiempo de emisión en formato UNIX
	}
	fmt.Printf("Esto es claimns: %+v\n", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey)) // TODO :: pasar por parámetro o atributo el algoritmo de encriptación
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, time.Now(), nil
}

// función para responder estructura correspondiente estandarizada
func responseJson(w http.ResponseWriter, detail string, data interface{}) {
	var response domainRequest.MockLoginResponse

	response.Detail = detail
	response.Data = data

	if detail == RESPONSE_OK {
		w.WriteHeader(http.StatusOK)
	}

	if detail == RESPONSE_ERROR {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(response)
	w.Header().Set("Content-Type", "application/json")
}
