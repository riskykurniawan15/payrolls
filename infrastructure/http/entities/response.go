package entities

type Response struct {
	Status           int         `json:"status"`
	Message          string      `json:"message"`
	Meta             interface{} `json:"meta,omitempty"`
	Data             interface{} `json:"data,omitempty"`
	Error            interface{} `json:"error,omitempty"`
	ValidationErrors interface{} `json:"validation_errors,omitempty"`
}

func ResponseFormater(Status int, data map[string]interface{}) Response {
	return Response{
		Status,
		GetResponseCodeMessage(Status),
		data["meta"],
		data["data"],
		data["error"],
		data["validation_errors"],
	}
}
