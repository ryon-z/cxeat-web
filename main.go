package main

import (
	"fmt"
	"os"
	envutil "yelloment-api/env_util"
	"yelloment-api/router"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	os.Setenv("TZ", envutil.GetGoDotEnvVariable("TIMEZONE"))
	if envutil.GetGoDotEnvVariable("MODE") == "OPERATION" {
		utils.SlackChannel = "operation"
		gin.SetMode(gin.ReleaseMode)
	} else {
		utils.SlackChannel = "debug"
	}
	port := envutil.GetGoDotEnvVariable("OPERATION_PORT")
	router := router.SetupRouter("operation_remote_test")
	// Temporary
	// models.CreateAllTableIfNotExists()

	router.Run(fmt.Sprintf(":%s", port))
}
