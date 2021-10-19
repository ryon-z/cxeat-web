package middleware

import (
	"strings"
	"yelloment-api/global"
	"yelloment-api/utils"

	"github.com/gin-gonic/gin"
)

// RenderWithParams : 파라미터를 포함하여 렌더
func RenderWithParams(statusCode int, pageID string, state gin.H) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		resultState := gin.H{}
		global.UpdateGlobalState(c)

		utils.CombineCustomStateGlobalState(c, &resultState)
		utils.CombineTwoGinH(&resultState, &state)

		pagePath := pageID
		if !strings.Contains(pageID, ".html") {
			pagePath = pageID + ".html"
		}
		c.HTML(statusCode, pagePath, resultState)
	}

	return gin.HandlerFunc(fn)
}

// RenderController : 파라미터를 포함하여 컨트롤러 렌더
func RenderController(controller gin.HandlerFunc) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		global.UpdateGlobalState(c)
		controller(c)
	}

	return gin.HandlerFunc(fn)
}
