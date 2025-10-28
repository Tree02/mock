package request

import "time"

type MockLoginRequest struct {
	HwId       *string `json:"hwid"`
	SucursalId *string `json:"sucursalid"`
	AgenteId   *string `json:"agenteid"`
	PuestoId   *string `json:"puestoid"`
	Pin        *string `json:"pin"`
}

type MockLoginResponse struct {
	Detail string      `json:"detail"`
	Data   interface{} `json:"data"`
}

// /getToken
// /refreshToken
type MockLoginData struct {
	Jwt   string    `json:"jwt"`
	Start time.Time `json:"start"`
}

// error
type MockLoginErrorData struct {
	Error string `json:"error"`
}
