package fastapi

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"reflect"
)

type EchoCtx struct {
	echo.Context
}

type Router struct {
	routesMap map[string]interface{}
}

func NewRouter() *Router {
	return &Router{
		routesMap: make(map[string]interface{}),
	}
}

func (r *Router) AddCall(path string, handler interface{}) {
	handlerType := reflect.TypeOf(handler)

	//if handlerType.NumIn() != 2 {
	//	panic("Wrong number of arguments")
	//}
	//if handlerType.NumOut() != 2 {
	//	panic("Wrong number of return values")
	//}

	//echoCtxType := reflect.TypeOf(&EchoCtx{})
	//if !handlerType.In(0).ConvertibleTo(echoCtxType) {
	//	panic("First argument should be *echo.Context!")
	//}
	fmt.Println(handlerType.In(0).Kind() == reflect.Struct)
	if handlerType.In(0).Kind() != reflect.Struct {
		panic("Second argument must be a struct")
	}

	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if !handlerType.Out(1).Implements(errorInterface) {
		panic("Second return value should be an error")
	}
	if handlerType.Out(0).Kind() != reflect.Struct {
		panic("First return value be a struct")
	}

	r.routesMap[path] = handler
}

func (r *Router) EchoHandler(c echo.Context) error {
	path := c.Param("path")
	log.Print(path)
	pathKey := "/" + path
	handlerFuncPtr, present := r.routesMap[pathKey]
	if !present {
		return fmt.Errorf("handler not found")
	}

	handlerType := reflect.TypeOf(handlerFuncPtr)
	inputType := handlerType.In(0)
	inputVal := reflect.New(inputType).Interface()
	//err := c.JSON(http.StatusOK, inputVal)
	//if err != nil {
	//	_ = c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	//	return
	//}

	toCall := reflect.ValueOf(handlerFuncPtr)
	outputVal := toCall.Call(
		[]reflect.Value{
			//reflect.ValueOf(c),
			reflect.ValueOf(inputVal).Elem(),
		},
	)

	returnedErr := outputVal[1].Interface()
	if returnedErr != nil || !outputVal[1].IsNil() {
		_ = c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": returnedErr})
		return fmt.Errorf("%v", returnedErr)
	}

	err := c.JSON(http.StatusOK, map[string]interface{}{"response": outputVal[0].Interface()})
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
