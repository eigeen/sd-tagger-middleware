package main

import (
	"ai-nsfw-detect/internal/config"
	"ai-nsfw-detect/internal/handler"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// 小工具，懒得分模块了，一个文件直接干
	config.LoadConfig()
	r := gin.Default()

	r.POST("/tagger/v1/interrogate", handler.Interrogate)

	r.Run(fmt.Sprintf("%s:%d",
		viper.GetString("server.host"),
		viper.GetInt("server.port")))
}
