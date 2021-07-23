# server
## 大致流程
server端监听端口6001,处理发送给/render的POST请求。

其中POST请求的Content-Type为`applcation/json`，必须包含字段`query_id`以及`formula`

server最后返回s3Url

注：如果想改动变量命名必须修改struct中的Tag值，Tag值与client发出的json信息中的key值必须完全对应