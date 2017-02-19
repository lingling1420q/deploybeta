package testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo"
)

func RequestJSON(handler echo.HandlerFunc, method string, url string, body interface{}) (*httptest.ResponseRecorder, map[string]string, error) {
	app := echo.New()

	jsonBody, err := json.Marshal(body)

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))

	if err != nil {
		panic(err)
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	res := httptest.NewRecorder()
	ctx := app.NewContext(req, res)
	err = handler(ctx)

	response := map[string]string{}

	if err == nil {
		err := json.Unmarshal(res.Body.Bytes(), &response)

		if err != nil {
			panic(err)
		}
	}

	return res, response, err
}
