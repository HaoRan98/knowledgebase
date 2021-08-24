package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gopkg.in/olahol/melody.v1"
	_ "knowledgebase/docs"
	"knowledgebase/middleware/cors"
	"knowledgebase/middleware/jwt"
	"knowledgebase/pkg/export"
	"knowledgebase/pkg/qrcode"
	"knowledgebase/pkg/upload"
	"knowledgebase/routers/api"
	v1 "knowledgebase/routers/api/v1"
	"log"
	"net/http"
	"time"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	//r.MaxMultipartMemory = int64(setting.AppSetting.FileMaxSize) << 20

	var mr = melody.New()
	mr.Config.MaxMessageSize = 40960 * 2
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.CORSMiddleware())

	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/", api.KlibIndex)
	r.GET("/klib", api.KlibIndex)
	r.Static("/css", "runtime/static/css")
	r.Static("/js", "runtime/static/js")
	r.Static("/img", "runtime/static/img")

	r.POST("/login", v1.Login)
	//r.POST("/login", v1.Login1)
	r.GET("/ws", v1.Websocket(mr))
	r.POST("/topic/imp", v1.ImpTopic)

	//代理转发智税平台登陆/获取路由表
	r.POST("/r_login", v1.Rlogin)
	r.GET("/r_route", v1.GetRoutes)

	// 智慧平台直接登陆
	r.POST("zhpt_login", v1.Zhpt_login)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		// 上传文件
		apiv1.POST("/file/upload", api.UploadFile)
		// 文件下载
		apiv1.StaticFS("/file/download", http.Dir(upload.GetFileFullPath()))
		// 删除附件
		apiv1.POST("/file/delete", api.DeleteFile)
		// 获取附件列表
		apiv1.POST("/files", api.GetFiles)

		// 发帖
		apiv1.POST("/topic/post", v1.PostTopic)
		// 修改帖子内容
		apiv1.POST("/topic/edit", v1.EditTopic)
		// 获取帖子详情
		apiv1.POST("/topic/detail", v1.GetTopic)
		// 获取帖子列表
		apiv1.GET("/topics", v1.GetTopics)
		// 置顶帖子
		apiv1.POST("/topic/top", v1.TopTopic)
		// 设置热门帖子
		apiv1.GET("/topic/hot/:id", v1.HotTopic)
		// 点赞帖子
		apiv1.GET("/topic/agree", v1.TopicAgree)
		// 取消点赞
		apiv1.GET("/topic/agree_cancel", v1.RemoveTopicAgree)
		// 删除帖子
		apiv1.GET("/topic/del/:id", v1.DelTopic)
		// 恢复帖子
		apiv1.POST("/topic/recovery", v1.RecTopic)
		// 查询删除帖子
		apiv1.POST("/topic/queryDel", v1.QueryDel)
		// 删除帖子 真
		apiv1.POST("/topic/delete", v1.DelTopic_zhen)
		// 查询机关单位帖子
		apiv1.POST("/topic/query", v1.QueryTopic)
		// 机关排名
		apiv1.POST("/topic/rank/dept", v1.Topic_dept_Rank)
		// 人员排名
		apiv1.POST("/topic/rank/user", v1.Topic_user_Rank)
		// 我的回帖和评论
		apiv1.POST("/repcom/lists", v1.GetRepCom)

		// 收藏帖子
		apiv1.GET("/favorite/add", v1.AddFavorite)
		// 取消收藏帖子
		apiv1.GET("/favorite/cancel", v1.CancelFavorite)
		// 获取收藏列表
		apiv1.GET("/favorites", v1.GetFavorites)
		// 删除收藏
		apiv1.GET("/favorite/del/:id", v1.DelFavorite)

		// 获取通知列表
		apiv1.GET("/notices", v1.GetNotices)
		// 删除通知
		apiv1.GET("/notice/del/:id", v1.DelNotice)

		// 获取分类列表
		apiv1.GET("/kinds", v1.GetKinds)
		// 删除分类
		apiv1.POST("/kind/delete", v1.DelKind)
		// 修改分类
		apiv1.POST("/kind/edit", v1.EditKinds)
		// 获取全部分类
		apiv1.GET("/kind/list", v1.GetKinds_zong)

		// 回帖
		apiv1.POST("/reply/post", v1.PostReply)
		// 修改回帖内容
		apiv1.POST("/reply/edit", v1.EditReply)
		// 获取回帖列表
		apiv1.GET("/replies", v1.GetReplies)
		// 采纳回帖
		apiv1.GET("/reply/accept/:id", v1.AcceptReply)
		// 点赞回帖
		apiv1.GET("/reply/agree/:id", v1.ReplyAgree)
		// 取消点赞回帖
		apiv1.GET("/reply/agree_cancel/:id", v1.RemoveReplyAgree)
		// 删除回帖
		apiv1.GET("/reply/del/:id", v1.DelReply)

		// 发布评论
		apiv1.POST("/comment/post", v1.PostComment)
		// 修改评论内容
		apiv1.POST("/comment/edit", v1.EditComment)
		// 获取评论列表
		apiv1.GET("/comments", v1.GetComments)
		// 点赞评论
		apiv1.GET("/comment/agree/:id", v1.CommentAgree)
		// 取消点赞评论
		apiv1.GET("/comment/agree_cancel/:id", v1.RemoveCommentAgree)
		// 删除评论
		apiv1.GET("/comment/del/:id", v1.DelComment)

		// 发布标签
		apiv1.POST("/label/post", v1.PostLabel)
		// 修改标签
		apiv1.POST("/label/edit", v1.EditLabel)
		// 获取标签列表
		apiv1.GET("/labels", v1.GetLabels)
		// 点赞标签
		apiv1.GET("/label/agree/", v1.LabelAgree)
		// 取消点赞标签
		apiv1.GET("/label/agree_cancel/", v1.RemoveLabelAgree)
		// 删除标签
		apiv1.GET("/label/del/", v1.DelLabel)

		// 获取人员列表
		apiv1.POST("/userinfo/", v1.GetUserList)
		// 创建团队
		apiv1.POST("/group/create", v1.CreateGroup)
		// 修改团队
		apiv1.POST("/group/edit", v1.EditGroup)
		// 获取团队列表
		apiv1.POST("/group/groups", v1.GetGroups)
		// 我创建的团队
		apiv1.POST("/group/mygroups", v1.MyGroups)
		// 删除团队
		apiv1.POST("/group/del", v1.DelGroup)
		// 查询团队
		apiv1.POST("/group/select", v1.SeleteGroup)

		// 添加成员
		apiv1.POST("/member/add", v1.AddMember)
		// 获取团队成员列表
		apiv1.POST("/member/list", v1.GetMembers)
		// 退出团队
		apiv1.POST("/member/dropout", v1.DropOut)

	}

	// 保证文本顺序输出
	go func() {
		time.Sleep(100 * time.Millisecond)
		// In order to ensure that the text order output can be deleted
		log.Println(`默认自动化文档地址:http://127.0.0.1:80/swagger/index.html`)
	}()
	return r
}
