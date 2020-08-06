package v1

import (
	"NULL/knowledgebase/models"
	"NULL/knowledgebase/pkg/app"
	"NULL/knowledgebase/pkg/e"
	"NULL/knowledgebase/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func AddFavorite(c *gin.Context) {
	appG := app.Gin{C: c}
	favorite := &models.Favorite{
		ID:      "FA-" + util.RandomString(29),
		TopicID: c.Query("id"),
		Account: util.GetLoginID("", c),
		Uptime:  time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.AddFavorite(favorite); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func CancelFavorite(c *gin.Context) {
	appG := app.Gin{C: c}
	topicId := c.Query("id")
	account := util.GetLoginID("", c)
	if err := models.CancelFavorite(topicId, account); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func GetFavorites(c *gin.Context) {
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
	favorites, err := models.GetFavorites(account, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(favorites) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS,
			map[string]interface{}{
				"list": favorites,
				"cnt":  models.GetFavoritesCnt(account),
			})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func DelFavorite(c *gin.Context) {
	var appG = app.Gin{C: c}
	if err := models.DelFavorite(c.Param("id")); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
