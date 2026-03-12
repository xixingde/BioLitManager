package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FileSizeLimit 创建文件大小限制中间件
// maxMB: 最大文件大小（MB）
func FileSizeLimit(maxMB int) gin.HandlerFunc {
	maxBytes := int64(maxMB * 1024 * 1024)

	return func(c *gin.Context) {
		// 检查 Content-Length 请求头
		if c.Request.ContentLength > maxBytes {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"code":    http.StatusRequestEntityTooLarge,
				"message": fmt.Sprintf("文件大小超过限制，最大允许 %dMB", maxMB),
			})
			return
		}

		// 设置 gin.Context 的 MaxBytesReader
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)

		c.Next()
	}
}
