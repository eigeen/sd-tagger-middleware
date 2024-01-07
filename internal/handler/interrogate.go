package handler

import (
	"ai-nsfw-detect/pkg/common"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"time"
)

type InterrogateData struct {
	Image     string  `json:"image,omitempty" form:"image"`
	Model     string  `json:"model,omitempty" form:"model"`
	Threshold float64 `json:"threshold,omitempty" form:"threshold"`
}

type ErrResp struct {
	Error  string `json:"error"`
	Detail string `json:"detail"`
	Body   string `json:"body"`
	Errors string `json:"errors"`
}

type SuccessResp struct {
	Caption map[string]interface{} `json:"caption"`
}

func Interrogate(ctx *gin.Context) {
	// 懒得分service层，直接写这里了
	var body InterrogateData
	err := ctx.ShouldBind(&body)
	if err != nil || body.Image == "" {
		ErrBadRequest(ctx, fmt.Sprintf("参数绑定错误 %s", err))
		return
	}

	// 覆盖参数
	if viper.GetString("override.model") != "" {
		body.Model = viper.GetString("override.model")
	}
	if viper.GetFloat64("override.threshold") != 0 {
		body.Threshold = viper.GetFloat64("override.model")
	}
	// 解码图片
	var img = make([]byte, base64.StdEncoding.DecodedLen(len(body.Image)))
	n, err := base64.StdEncoding.Decode(img, []byte(body.Image))
	if err != nil {
		ErrImageEncoding(ctx, fmt.Sprintf("图片base64解码失败 %s", err))
		return
	}
	img = img[:n]

	// 图片写入临时文件
	fileName := uuid.New().String()
	filePath := fmt.Sprintf("images/%s.jpg", fileName)
	exportFilePath := fmt.Sprintf("images/%s_exp.webp", fileName)
	f, err := common.CreateNestedFile(filePath)
	defer f.Close()
	if err != nil {
		ErrInternalServer(ctx, fmt.Sprintf("创建临时文件失败 %s", err))
		return
	}

	_, err = f.Write(img)
	if err != nil {
		ErrInternalServer(ctx, fmt.Sprintf("写入临时文件失败 %s", err))
		return
	}

	// 压缩图像
	err = common.ImageCompress(filePath, exportFilePath)
	if err != nil {
		ErrInternalServer(ctx, fmt.Sprintf("压缩图像出错 %s", err))
		return
	}

	// 读取压缩后图像并编码
	fexp, err := os.Open(exportFilePath)
	defer fexp.Close()
	if err != nil {
		ErrInternalServer(ctx, fmt.Sprintf("读取压缩后图像出错 %s", err))
		return
	}

	expImg, err := io.ReadAll(fexp)
	if err != nil {
		ErrInternalServer(ctx, fmt.Sprintf("读取压缩后图像出错 %s", err))
		return
	}

	expImgStr := base64.StdEncoding.EncodeToString(expImg)
	body.Image = expImgStr

	// 请求后端 并将返回值原样递送回客户端
	jsonBody, err := json.Marshal(&body)
	if err != nil {
		ErrInternalServer(ctx, fmt.Sprintf("序列化参数失败 %s", err))
		return
	}

	client := http.Client{Timeout: time.Second * 30}
	backendUrl := viper.GetString("backend.base_url") + "/tagger/v1/interrogate"
	resp, err := client.Post(backendUrl,
		"application/json",
		bytes.NewReader(jsonBody))
	if err != nil {
		ErrInternalServer(ctx, fmt.Sprintf("请求后端失败 %s", err))
		return
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ErrInternalServer(ctx, fmt.Sprintf("读取后端返回值失败 %s", err))
		return
	}

	var respBodyMap = make(map[string]any)
	err = json.Unmarshal(respBody, &respBodyMap)
	if err != nil {
		ErrInternalServer(ctx, fmt.Sprintf("读取后端返回值失败 %s", err))
		return
	}

	ctx.JSON(resp.StatusCode, respBodyMap)
}
