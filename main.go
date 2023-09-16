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

	routes.UserRoutes(r)
	r.Run()
}
