package main

import (
	"theatreManagementApp/config"
	"theatreManagementApp/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnvVariables()
	config.ConnectToDB()
}

func main() {
	r := gin.Default()

	// Serve HTML templates from the "templates" directory
	r.LoadHTMLGlob("views/*.html")

	routes.UserRoutes(r)
	routes.AdminRoutes(r)
	routes.ManagerRoutes(r)
	r.Run()
}
