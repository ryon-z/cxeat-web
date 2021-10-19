package viewcontrollers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Logout : 로그아웃
func Logout(c *gin.Context) {
	cueatMainDomain := ".cueat.kr"

	candidateDomains := []string{
		"",
		cueatMainDomain,
		"test" + cueatMainDomain,
		"test-api" + cueatMainDomain,
		"test-admin" + cueatMainDomain,
		"m.test" + cueatMainDomain,
	}

	for _, domain := range candidateDomains {
		c.SetCookie("jwt", "", -1, "/", domain, false, true)
		c.SetCookie("expire", "", -1, "/", domain, false, true)
		c.SetCookie("pageAfterLogin", "", -1, "/", domain, false, true)
	}

	c.Set("loggedIn", "no")

	homePath := "/web/main"
	c.Redirect(http.StatusTemporaryRedirect, homePath)
}
