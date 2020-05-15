package routers

import (
	"NULL/knowledgebase/middleware/cors"
	"NULL/knowledgebase/middleware/jwt"
	"NULL/knowledgebase/pkg/export"
	"NULL/knowledgebase/pkg/qrcode"
	"NULL/knowledgebase/pkg/upload"
	"NULL/knowledgebase/routers/api"
	v1 "NULL/knowledgebase/routers/api/v1"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"net/http"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	var mr = melody.New()
	mr.Config.MaxMessageSize = 40960 * 2
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.CORSMiddleware())

	r.GET("/", v1.Home)
	r.GET("/klib", v1.Home)
	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))
	r.Static("/css", "runtime/static/css")
	r.Static("/js", "runtime/static/js")
	r.Static("/img", "runtime/static/img")

	r.POST("/login", v1.Login)
	r.GET("/ws", v1.Websocket(mr))
	r.POST("/topic/imp", v1.ImpTopic)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		//上传文件
		apiv1.POST("/file/upload", api.UploadFile)
		//文件下载
		apiv1.StaticFS("/file/download", http.Dir(upload.GetFileFullPath()))

		//发帖
		apiv1.POST("/topic/post", v1.PostTopic)
		//修改帖子内容
		apiv1.POST("/topic/edit", v1.EditTopic)
		//获取帖子详情
		apiv1.GET("/topic/detail/:id", v1.GetTopic)
		//获取帖子列表
		apiv1.GET("/topics", v1.GetTopics)
		//置顶帖子
		apiv1.GET("/topic/top/:id", v1.TopTopic)
		//设置热门帖子
		apiv1.GET("/topic/hot/:id", v1.HotTopic)
		//删除帖子
		apiv1.GET("/topic/del/:id", v1.DelTopic)

		//收藏帖子
		apiv1.GET("/favorite/add", v1.AddFavorite)
		//获取收藏列表
		apiv1.GET("/favorites", v1.GetFavorites)
		//删除收藏
		apiv1.GET("/favorite/del/:id", v1.DelFavorite)

		//获取通知列表
		apiv1.GET("/notices", v1.GetNotices)
		//删除通知
		apiv1.GET("/notice/del/:id", v1.DelNotice)

		//获取分类列表
		apiv1.GET("/kinds", v1.GetKinds)
		//删除分类
		apiv1.GET("/kind/del/:id", v1.DelKind)

		//回帖
		apiv1.POST("/reply/post", v1.PostReply)
		//修改回帖内容
		apiv1.POST("/reply/edit", v1.EditReply)
		//获取回帖列表
		apiv1.GET("/replies", v1.GetReplies)
		//采纳回帖
		apiv1.GET("/reply/accept/:id", v1.AcceptReply)
		//点赞回帖
		apiv1.GET("/reply/agree/:id", v1.ReplyAgree)
		//删除回帖
		apiv1.GET("/reply/del/:id", v1.DelReply)

		//发布评论
		apiv1.POST("/comment/post", v1.PostComment)
		//修改评论内容
		apiv1.POST("/comment/edit", v1.EditComment)
		//获取回帖列表
		apiv1.GET("/comments", v1.GetComments)
		//点赞评论
		apiv1.GET("/comment/agree/:id", v1.CommentAgree)
		//删除评论
		apiv1.GET("/comment/del/:id", v1.DelComment)
	}
	return r
}
