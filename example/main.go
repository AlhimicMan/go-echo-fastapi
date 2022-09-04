package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-fastapi"
	"net/http"
)

type EchoInput struct {
	Phrase string `json:"phrase"`
}

type EchoOutput struct {
	OriginalInput EchoInput `json:"original_input"`
}

func EchoHandler(in EchoInput) (out EchoOutput, err error) {
	out.OriginalInput = in
	return
}

func GenerateScheme() {
	myRouter := fastapi.NewRouter()
	myRouter.AddCall("/echo", EchoHandler)

	swagger := myRouter.EmitOpenAPIDefinition()
	swagger.Info.Title = "My awesome API"
	jsonBytes, _ := json.MarshalIndent(swagger, "", "    ")
	fmt.Println(string(jsonBytes))
}

func main() {
	GenerateScheme()

	r := echo.New()

	myRouter := fastapi.NewRouter()
	myRouter.AddCall("/echo", EchoHandler)

	r.POST("/api/:path", myRouter.EchoHandler) // must have *path parameter

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:8034"),
		Handler: r,
	}

	_ = httpServer.ListenAndServe()
}
