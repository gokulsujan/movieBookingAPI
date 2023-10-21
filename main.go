package main

import (
	"text/template"
	"theatreManagementApp/config"
	"theatreManagementApp/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnvVariables()
	config.ConnectToDB()
}

func main() {
	// Define a custom "add" function for serial numbers
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"eq": func(a, b interface{}) bool {
			return a == b
		},
	}

	r := gin.Default()

	// Register the custom function with the template engine
	r.SetFuncMap(funcMap)

	// Serve HTML templates from the "templates" directory
	r.LoadHTMLGlob("views/*.html")

	routes.UserRoutes(r)
	routes.AdminRoutes(r)
	routes.ManagerRoutes(r)
	r.Run()
}
