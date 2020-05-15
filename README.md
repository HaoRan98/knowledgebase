# KnowledgeBase

## Es API
```
Get document by ID:
  curl -X GET "localhost:9200/index/_source/id?pretty"
Get All documents:
  curl -X GET "localhost:9200/_search?pretty"
Delete document:
  curl -X DELETE "localhost:9200/index/_doc/id?routing=kimchy&pretty"
```
  
## WebSocket

 + websocket统一消息格式
 ```
 {
    "ms_type":"",
    "message":string,   //客户端推送服务端，解析数据
    "data":object,      //服务端推送客户端，解析数据
    "token":string,  
 }
 ```

 + ms_type   类型及data对象
 ```
 鉴权：                
   {"ms_type":"check","data":"AUTH CHECK TOKEN FAIL !!!"  }
 在线人数：              
   {"ms_type":"online","data":number}
 获取search结果：        
   {"ms_type":"esrecord","data":失败返回string，成功返回对象数组  }  
 获取topic列表，返回与接口相同:     
   {"ms_type":"topic","data":map}
 向用户广播消息通知:                 
   {"ms_type":"notice","data":Notice表对象}
 向用户广播各类统计信息:           
   {"ms_type":"count","data":map}  

```
  