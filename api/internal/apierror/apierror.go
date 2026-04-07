package apierror

import "encoding/json"

type APIError struct {
	Status    int    `json:"status"`
	ErrorCode string `json:"error_code"`
	Title     string `json:"title"`
	Details   string `json:"details"`
}

func (e *APIError) JSON() []byte {
	b, _ := json.Marshal(e)
	return b
}