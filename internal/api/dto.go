package api

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type CreateAPIRequest struct {
	Name    string `json:"name" binding:"required"`
	Tag     string `json:"tag"`
	APIURL  string `json:"api_url" binding:"required"`
	Token   string `json:"token" binding:"required"`
	Model   string `json:"model" binding:"required"`
	Enabled bool   `json:"enabled"`
}

type UpdateAPIRequest struct {
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	APIURL  string `json:"api_url"`
	Token   string `json:"token"`
	Model   string `json:"model"`
	Enabled *bool  `json:"enabled"`
}

func Success(data interface{}) Response {
	return Response{Code: 0, Message: "ok", Data: data}
}

func Error(message string) Response {
	return Response{Code: 1, Message: message, Data: nil}
}
