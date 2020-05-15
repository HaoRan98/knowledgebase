package v1

import (
	"NULL/knowledgebase/models"
	"NULL/knowledgebase/pkg/app"
	"NULL/knowledgebase/pkg/e"
	"NULL/knowledgebase/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetNotices(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		pageSize int
		pageNo   int
	)
	account := util.GetLoginID("", c)
	if c.Query("pageNo") == "" {
		pageNo = 1
	} else {
		pageNo, _ = strconv.Atoi(c.Query("pageNo"))
	}
	if c.Query("pageSize") == "" {
		pageSize = 100
	} else {
		pageSize, _ = strconv.Atoi(c.Query("pageSize"))
	}
	favorites, err := models.GetNotices(account, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(favorites) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS,
			map[string]interface{}{
				"list": favorites,
				"cnt":  models.GetNoticesCnt(account),
			})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func DelNotice(c *gin.Context) {
	var appG = app.Gin{C: c}
	if err := models.DelNotice(c.Param("id")); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
