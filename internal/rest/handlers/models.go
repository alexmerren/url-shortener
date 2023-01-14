package handlers

type addURLRequestBody struct {
	Url string `json:"url"`
}

type addURLResponseBody struct {
	ID string `json:"id"`
}

type getURLResponseBody struct {
	Url string `json:"url"`
}

type errorResponse struct {
	StatusCode  int    `json:"status_code"`
	ErrorString string `json:"error_string"`
}
