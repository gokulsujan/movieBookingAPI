package main

import "theatreManagementApp/config"

func init() {
	config.LoadEnvVariables()
	config.ConnectToDB()
}

func main() {

}
