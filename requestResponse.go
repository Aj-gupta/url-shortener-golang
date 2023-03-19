package urlshortner

type GenericResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
