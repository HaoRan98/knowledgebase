package v1

import (
	"NULL/knowledgebase/models"
	"NULL/knowledgebase/pkg/app"
	"NULL/knowledgebase/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetKinds(c *gin.Context) {
	appG := app.Gin{C: c}
	kinds, err := models.GetKinds()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(kinds) > 0 {
		mc := make([]string, len(kinds))
		for _, kind := range kinds {
			mc = append(mc, kind.Mc)
		}
		appG.Response(http.StatusOK, e.SUCCESS, mc)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func DelKind(c *gin.Context) {
	var appG = app.Gin{C: c}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if err := models.DelKind(uint(id)); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
