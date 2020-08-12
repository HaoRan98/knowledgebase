package api

import (
	"NULL/knowledgebase/pkg/util"
	"github.com/gin-gonic/gin"
	"html/template"
)

func KlibIndex(c *gin.Context) {
	if c.Request.Method == "GET" {
		t, _ := template.ParseFiles("runtime/static/index.html")
		util.ShowError("template parseFiles err", t.Execute(c.Writer, nil))
	}
}

func JkxmIndex(c *gin.Context) {
	if c.Request.Method == "GET" {
		t, _ := template.ParseFiles("runtime/static_jkxm/index.html")
		util.ShowError("template parseFiles err", t.Execute(c.Writer, nil))
	}
}

func KfqZmq1(c *gin.Context) {
	if c.Request.Method == "GET" {
		t, _ := template.ParseFiles("runtime/static_kfqzmq1/index.html")
		util.ShowError("template parseFiles err", t.Execute(c.Writer, nil))
	}
}

func KfqZmq2(c *gin.Context) {
	if c.Request.Method == "GET" {
		t, _ := template.ParseFiles("runtime/static_kfqzmq2/index.html")
		util.ShowError("template parseFiles err", t.Execute(c.Writer, nil))
	}
}
