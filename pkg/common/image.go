package common

import ffmpeg "github.com/u2takey/ffmpeg-go"

// ImageCompress 转码并压缩图像
//
// 编码格式取决于提供的输出文件后缀名，暂不支持显式指定-c:v
// 分辨率设置为宽512x，高度自适应
func ImageCompress(input, output string) error {
	err := ffmpeg.Input(input).
		Output(output, ffmpeg.KwArgs{"vf": "scale=512:-1"}).
		OverWriteOutput().
		ErrorToStdOut().
		Run()

	return err
}
