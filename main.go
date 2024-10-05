package main

import (
	"api-dot/router"
	"api-dot/utils"
)

func main() {
	route := router.SetupRouter(utils.GetConfig("GIN_MODE"))

	route.Run(":" + utils.GetConfig("PORT"))
}