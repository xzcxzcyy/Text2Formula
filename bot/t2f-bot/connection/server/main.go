package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type PictureRender struct {
	QueryID string `json:"query_id"`
	Formula string `json:"formula"`
}

var (
	Url = "https://www.baidu.com"
)


func InitServer() {
	r := gin.Default()

	r.POST("/render", func(c *gin.Context) {
		pictureRender := PictureRender{}

		// get the info
		err := c.BindJSON(&pictureRender)
		if err != nil {
			log.Println(err)
			return
		}
		// 验证是否拿到了正确的信息
		//log.Printf("queryid:%+v\nforumla:%+v\n", pictureRender.QueryID, pictureRender.Formula)

		// 调用函数，然后最后会得到s3的url
		// ...

		// return directly
		c.JSON(http.StatusOK, gin.H{
			"s3Url": Url,
		})
	})
	r.Run(":6001")
}
