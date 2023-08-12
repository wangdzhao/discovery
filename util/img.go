package util

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func SendImageToWeChatRobot(imagePath string) (string, string, error) {
	// 读取图片文件
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return "", "", err
	}
	defer imageFile.Close()

	imageData, err := ioutil.ReadAll(imageFile)
	if err != nil {
		return "", "", err
	}

	// 计算图片的md5值
	md5Sum := md5.Sum(imageData)
	md5Param := fmt.Sprintf("%x", md5Sum)

	// 将图片数据转换为base64编码
	base64Data := base64.StdEncoding.EncodeToString(imageData)

	return md5Param, base64Data, nil
}

// 获取整个浏览器窗口的截图（全屏）
// 这将模拟浏览器操作设置。
func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(10 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// 得到布局页面
			_, _, _, _, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}
			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))
			// 浏览器视窗设置模拟
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// 捕捉屏幕截图
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      730,
					Y:      310,
					Width:  450,
					Height: 450,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}

func SaveCodeImg(url string) {
	// 禁用chrome headless
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoDefaultBrowserCheck, //不检查默认浏览器
		chromedp.Flag("headless", true),
		chromedp.Flag("blink-settings", "imagesEnabled=true"), //开启图像界面,重点是开启这个
		chromedp.Flag("ignore-certificate-errors", true),      //忽略错误
		chromedp.Flag("disable-web-security", true),           //禁用网络安全标志
		chromedp.Flag("disable-extensions", true),             //开启插件支持
		chromedp.Flag("disable-default-apps", true),
		chromedp.WindowSize(1920, 1080),    // 设置浏览器分辨率（窗口大小）
		chromedp.Flag("disable-gpu", true), //开启gpu渲染
		chromedp.Flag("hide-scrollbars", true),
		chromedp.Flag("mute-audio", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.NoFirstRun, //设置网站不是首次运行
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36"), //设置UserAgent
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// 创建上下文实例
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(LogInfo),
	)
	defer cancel()

	// 创建超时上下文
	ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	//导航到目标页面，等待一个元素，捕捉元素的截图
	var buf []byte
	if err := chromedp.Run(ctx, fullScreenshot(url, 100, &buf)); err != nil {
		LogError("截图失败: %v", err)
		return
	}
	if err := ioutil.WriteFile("code.png", buf, 0644); err != nil {
		LogError("图片写入失败: %v", err)
		return
	}
	LogInfo("图片写入完成")
}
