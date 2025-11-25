package router

import "github.com/gin-gonic/gin"

// Test
// @Tags  接口分组
// @Summary 测试接口
// @Router    /api/v1/test [get]
func Test(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":    200,
		"message": "test",
		"data":    gin.H{"name": "test"},
	})
}
