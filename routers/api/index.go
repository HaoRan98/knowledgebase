package api

import (
	"NULL/knowledgebase/pkg/util"
	"github.com/gin-gonic/gin"
	"html/template"
)

func Home(c *gin.Context) {
	if c.Request.Method == "GET" {
		t, _ := template.ParseFiles("runtime/static/index.html")
		util.ShowError("template parseFiles err", t.Execute(c.Writer, nil))
	}
}
