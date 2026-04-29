// SPDX-Licence-Identifier: EUPL-1.2

package api

import (
	. "dappco.re/go"
	"github.com/gin-gonic/gin"
)

func ExampleNewProvider() {
	Println(NewProvider().Name())
	// Output: html
}

func ExampleHTMLProvider_Name() {
	provider := &HTMLProvider{}
	Println(provider.Name())
	// Output: html
}

func ExampleHTMLProvider_BasePath() {
	provider := &HTMLProvider{}
	Println(provider.BasePath())
	// Output: /v1/html
}

func ExampleHTMLProvider_RegisterRoutes() {
	gin.SetMode(gin.TestMode)
	provider := NewProvider()
	router := gin.New()
	provider.RegisterRoutes(router.Group(provider.BasePath()))
	Println(len(router.Routes()))
	// Output: 2
}

func ExampleHTMLProvider_Describe() {
	routes := NewProvider().Describe()
	Println(len(routes), routes[0].Path)
	// Output: 2 /render
}
