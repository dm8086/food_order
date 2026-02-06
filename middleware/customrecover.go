package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// 1. 打印错误信息和调用栈
				fmt.Printf("Panic recovered: %v\n", r)
				// fmt.Printf("Stack trace: %s\n", debug.Stack())
			}
		}()
		c.Next()
	}
}
