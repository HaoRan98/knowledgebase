package v1

import (
	"NULL/knowledgebase/models"
	"NULL/knowledgebase/pkg/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"html/template"
	"log"
	"sync"
)

type WsMsg struct {
	MsType  string      `json:"ms_type"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Token   string      `json:"token"`
}

//连接列表
//map[*melody.Session]*models.SysUser
var Clients = sync.Map{}

func Websocket(mr *melody.Melody) gin.HandlerFunc {
	mr.HandleDisconnect(func(session *melody.Session) {
		if u, ok := Clients.Load(session); ok {
			log.Printf("###[%v] client ws disconnect", u)
			Clients.Delete(session)
			BroadCastOnline()
		}
	})
	mr.HandleError(func(session *melody.Session, e error) {
		if u, ok := Clients.Load(session); ok {
			log.Printf("###[%v] client ws error:  err: %v", u, e)
			Clients.Delete(session)
			BroadCastOnline()
		}
	})
	mr.HandleMessage(func(session *melody.Session, msg []byte) {
		var wsMsg = WsMsg{}
		log.Printf("rec: %v", string(msg))
		util.ShowError("", json.Unmarshal(msg, &wsMsg))
		if util.GetLoginID(wsMsg.Token, nil) == "" {
			msg := WsMsg{
				MsType: "check",
				Data:   "AUTH CHECK TOKEN FAIL !!!",
			}
			msgJsons, err := json.Marshal(msg)
			if err != nil {
				log.Printf("search marshal json.\n ERR:%v", err)
				return
			}
			Send(session, msgJsons)
			return
		}
		switch wsMsg.MsType {
		case "hello":
			var u = &models.SysUser{}
			util.ShowError("", util.FormJson(wsMsg.Message, &u))
			log.Println(u)
			Clients.Store(session, u)
			BroadCastOnline()
			BroadCastTopic()
			BroadCastCount()
		case "search":
			esMap := make(map[string]string)
			util.ShowError("rec err", util.FormJson(wsMsg.Message, &esMap))
			log.Println(esMap)

			//搜索量+1
			var lock sync.RWMutex
			lock.Lock()
			models.Info.Browse++
			lock.Unlock()

			var index string
			if esMap["index"] == "" {
				index = "_all"
			} else {
				index = esMap["index"]
			}
			msg := WsMsg{
				MsType: "esrecord",
			}
			resp, errMsg := models.EsSearch(index, esMap["key"])
			if errMsg != "" {
				log.Println("EsDocument err:", errMsg)
				msg.Data = errMsg
			} else {
				//log.Println("EsDocument:", resp)
				msg.Data = resp
			}
			msgJsons, err := json.Marshal(msg)
			if err != nil {
				log.Printf("search marshal json.\n ERR:%v", err)
				return
			}
			Send(session, msgJsons)
		case "bye":
			if u, ok := Clients.Load(session); ok {
				log.Printf("###client ws disconnect: %v", u)
				Clients.Delete(session)
				BroadCastOnline()
			}
		}
	})
	return func(c *gin.Context) {
		util.ShowError("websocket handle err", mr.HandleRequest(c.Writer, c.Request))
		c.Next()
	}
}

func Send(session *melody.Session, m []byte) {
	err := session.Write(m)
	if err != nil {
		log.Println(err)
	}
}

// 新连接或断开后,向所有人广播在线人数
func BroadCastOnline() {
	var onlineNum int
	Clients.Range(func(s, u interface{}) bool {
		onlineNum++
		return true
	})
	msg := WsMsg{
		MsType: "online",
		Data:   onlineNum,
	}
	msgJsons, err := json.Marshal(msg)
	if err != nil {
		log.Printf("broadcast online fail at marshal json.\n ERR:%v", err)
		return
	}
	Clients.Range(func(s, u interface{}) bool {
		err := s.(*melody.Session).Write(msgJsons)
		if err != nil {
			log.Printf("send online msg err\n ERR:%v", err)
		}
		return true
	})
}

// 发帖或修改帖子后,向所有人广播新帖子列表
func BroadCastTopic() {
	topics, err := models.GetTopics("", "", 1, 10)
	if err != nil {
		log.Printf("send topic msg err\n ERR:%v", err)
		return
	}
	tpResps := make([]TpResp, 0)
	for _, tp := range topics {
		flag := models.IsAgreed(tp.ID)
		tpResps = append(tpResps, TpResp{tp, flag})
	}
	resp := map[string]interface{}{
		"list": tpResps,
		"cnt":  models.GetTopicsCnt(""),
	}
	msg := WsMsg{
		MsType: "topic",
		Data:   resp,
	}
	msgJsons, err := json.Marshal(msg)
	if err != nil {
		log.Printf("broadcast topic fail at marshal json.\n ERR:%v", err)
		return
	}
	Clients.Range(func(s, u interface{}) bool {
		err := s.(*melody.Session).Write(msgJsons)
		if err != nil {
			log.Printf("send topic msg err\n ERR:%v", err)
		}
		return true
	})
}

// 回帖后,向发帖人、收藏人或回帖人广播通知,在线广播;不在线,写入通知表
func BroadCastReply(notice *models.Notice) {
	msg := WsMsg{
		MsType: "notice",
		Data:   notice,
	}
	msgJsons, err := json.Marshal(msg)
	if err != nil {
		log.Printf("broadcast notice fail at marshal json.\n ERR:%v", err)
		return
	}

	var online bool
	Clients.Range(func(s, u interface{}) bool {
		if u.(*models.SysUser).UserAccount == notice.Account {
			online = true
			err := s.(*melody.Session).Write(msgJsons)
			if err != nil {
				log.Printf("send notice msg err\n ERR:%v", err)
			}
		}
		return true
	})

	if !online {
		err = models.AddNotice(notice)
		if err != nil {
			log.Printf("add notice to db err\n ERR:%v", err)
		}
	}
}

// 发布或修改topic、reply、comment,浏览量增加后,向所有人广播新统计信息
func BroadCastCount() {
	var wordCnt int
	var lock sync.RWMutex
	defer lock.Unlock()
	tName := []string{"reply", "comment"}
	for _, tname := range tName {
		wordCnt += models.GetWordCnt(tname)
	}
	wordCnt += models.GetTopicWordCnt()
	lock.Lock()
	//统计信息
	count := map[string]int{
		"topicNum":   models.GetTopicsCnt(""),    //发帖数
		"replyNum":   models.GetRepliesCnt(""),   //回帖数
		"commentNum": models.GetCommentsCnt(""),  //评论数
		"browseCnt":  models.GetTopicBrowseCnt(), //总浏览量
		"searchCnt":  models.Info.Browse,         //搜索量
		"wordCnt":    wordCnt,                    //字数统计
	}
	msg := WsMsg{
		MsType: "count",
		Data:   count,
	}
	msgJsons, err := json.Marshal(msg)
	if err != nil {
		log.Printf("broadcast count info fail at marshal json.\n ERR:%v", err)
		return
	}
	Clients.Range(func(s, u interface{}) bool {
		err := s.(*melody.Session).Write(msgJsons)
		if err != nil {
			log.Printf("send count info msg err\n ERR:%v", err)
		}
		return true
	})
}

func Home(c *gin.Context) {
	if c.Request.Method == "GET" {
		t, _ := template.ParseFiles("runtime/static/index.html")
		util.ShowError("template parseFiles err", t.Execute(c.Writer, nil))
	}
}
