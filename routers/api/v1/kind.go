package v1

import (
	"github.com/gin-gonic/gin"
	"knowledgebase/models"
	"knowledgebase/pkg/app"
	"knowledgebase/pkg/e"
	"log"
	"net/http"
	"strconv"
)

type Kind struct {
	ID string `json:"id"`
	MC string `json:"mc"`
}

func GetKinds(c *gin.Context) {
	appG := app.Gin{C: c}
	kinds, err := models.GetKinds()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(kinds) > 0 {
		mc := make([]string, 0)
		for _, kind := range kinds {
			if kind.Mc != "" {
				mc = append(mc, kind.Mc)
			}
		}
		appG.Response(http.StatusOK, e.SUCCESS, mc)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func EditKinds(c *gin.Context) {

	var (
		appG = app.Gin{C: c}
		form Kind
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	if !models.IsKindExist(form.MC) {
		id, err := strconv.Atoi(form.ID)
		if err != nil {
			log.Println("kind id转int失败")
			appG.Response(httpCode, errCode, nil)
			return
		}

		kind := map[string]interface{}{
			"id": id,
			"mc": form.MC,
		}

		kindmc, err := models.Getkind(uint(id))
		if err != nil {
			log.Println(err)
			appG.Response(http.StatusInternalServerError, e.ERROR, "查询分类失败")
			return
		}

		if err := models.EditKind(kind); err != nil {
			log.Println(err)
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
			return
		}

		if err := models.UpDateKind(kindmc.Mc, form.MC); err != nil {
			log.Println(err)
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
			return
		}

		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	}

	appG.Response(http.StatusOK, e.ERROR, "该分类已存在")

}

func DelKind(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form Kind
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	id, err := strconv.Atoi(form.ID)
	if err != nil {
		log.Println("kind id转int失败")
		appG.Response(httpCode, errCode, nil)
		return
	}

	err, ok := models.GetKind(uint(id))
	if err != nil {
		log.Println(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	if !ok {
		if err := models.DelKind(uint(id)); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, "删除成功")
	} else {
		appG.Response(http.StatusOK, e.ERROR, "该分类已被使用，无法删除")
		return
	}

}

func GetKinds_zong(c *gin.Context) {
	appG := app.Gin{C: c}
	kinds, err := models.GetKinds()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(kinds) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, kinds)
		return
	}
	appG.Response(http.StatusOK, e.ERROR, "分类获取为空")
}
