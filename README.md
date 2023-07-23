<p align="center">
<a href="https://pkg.go.dev/github.com/adam-qiang/useragent"><img src="https://pkg.go.dev/badge/github.com/adam-qiang/useragent.svg" alt="Go Reference"></a>
<a href="https://en.wikipedia.org/wiki/MIT_License" rel="nofollow"><img alt="MIT" src="https://img.shields.io/badge/license-MIT-blue.svg" style="max-width:100%;"></a>
</p>

---

# useragent

客户端请求代理信息解析

**注：**

- 由于是根据User-Agent字符的规则匹配（可能因浏览器、设备和操作系统的不同而有所变化），所以可能会存在识别错误或识别不全的问题
- 由于对于Windows11微软更改了规范所以无法识别，结果将和Windows10一样，相关文档：
  https://learn.microsoft.com/zh-cn/microsoft-edge/web-platform/how-to-detect-win11
- 由于实现方案是进行规则匹配故而有可能会有些兼容不到的情况，如有发现可随时提交问题
- 关于机器人目前只判断了Google和Bing，后续会建立机器人库（即枚举）进行匹配

## 安装

```go
 go get github.com/adam-qiang/useragent@latest
```

## Demo

```go
package main

import (
	"fmt"
	"github.com/adam-qiang/useragent"
)

func main() {
	userAgent := "Mozilla/5.0 (Linux; Android 12; HarmonyOS; NCO-AL00; HMSCore 6.11.0.302) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.88 HuaweiBrowser/14.0.0.322 Mobile Safari/537.36"
	ua := useragent.New(userAgent)

	fmt.Println("Is Mobile Device：" + fmt.Sprint("", ua.IsMobile))
	fmt.Println("Device Name：" + ua.Mobile)
	fmt.Println("Platform：" + ua.Os.Platform)
	fmt.Println("System Name：" + ua.Os.Name)
	fmt.Println("System Version：" + ua.Os.Version)

	fmt.Println("Browser Id:" + ua.Browser.Id)
	fmt.Println("Browser Id Version:" + ua.Browser.IdVersion)
	fmt.Println("Browser：" + ua.Browser.Name)
	fmt.Println("Browser Version：" + ua.Browser.Version)
	fmt.Println("Engine：" + ua.Browser.Engine)
	fmt.Println("Engine Version：" + ua.Browser.EngineVersion)
}
```