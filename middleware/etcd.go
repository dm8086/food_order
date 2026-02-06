package middleware

import (
	"encoding/base64"
	"net/http"
	"order_food/global"
	"order_food/model/common/response"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func InternalCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.Path, "internal/") {
			// 需要校验服务的uuid
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusOK, gin.H{"code": response.ERROR, "msg": "内部调用接口参数缺失"})
				c.Abort()
				return
			}
			// 解析里面的值
			authHeaderByte, _ := base64.StdEncoding.DecodeString(authHeader)
			authors := strings.Split(string(authHeaderByte), ";")
			if len(authors) != 3 {
				c.JSON(http.StatusOK, gin.H{"code": response.ERROR, "msg": "内部调用接口参解析错误"})
				c.Abort()
				return
			}

			if authors[0] != global.ServUuid {
				c.JSON(http.StatusOK, gin.H{"code": response.ERROR, "msg": "内部调用接口解密错误"})
				c.Abort()
				return
			}
			reqTime, _ := strconv.ParseInt(authors[2], 10, 64)
			if time.Now().Local().Unix()-reqTime > 15 {
				c.JSON(http.StatusOK, gin.H{"code": response.ERROR, "msg": "内部调用接口时间大于15秒"})
				c.Abort()
				return
			}

			c.Next()
			return
		}
		c.Next()
	}
}
