/**
 * Created by goland.
 * User: adam_wang
 * Date: 2023-07-14 00:37:31
 */

package useragent

import (
	"regexp"
	"strings"
)

// Browser
// @Description: 浏览器结构体
type Browser struct {
	Id            string //浏览器标识
	IdVersion     string //浏览器标识版本
	Name          string
	Version       string
	Engine        string //浏览器内核（即使用的渲染引擎）
	EngineVersion string //内核版本
}

// 解析浏览器信息
// @receiver u *UserAgent
// @param part []Part
func (u *UserAgent) analysisBrowser(part []Part) {
	partLen := len(part)
	isWechat := strings.HasPrefix(u.Ua, "MicroMessenger")

	if partLen > 1 {
		engine := part[1]
		u.Browser.Engine = engine.Name
		u.Browser.EngineVersion = engine.Version

		//TODO：需要增加一个判断是否是机器人的正则逻辑，暂时先判断goole和bing的机器人
		for _, other := range engine.Other {
			if strings.HasPrefix(other, "Googlebot") || strings.HasPrefix(other, "bingbot") {
				u.IsRobot = true
				break
			}
		}

		//如果是手机端则需要查看是否是在微信端或者是否是定制浏览器
		if u.IsMobile {
			for _, partInfo := range part {
				if (isWechat && partInfo.Name == "MicroMessenger") || strings.Contains(partInfo.Name, "Browser") {
					u.Browser.Name = partInfo.Name
					u.Browser.Version = partInfo.Version

					return
				}
			}
		}

		if partLen > 2 {
			partIndex := 2
			//解决个别浏览器内核后面的版本号没有
			if part[2].Version == "" && partLen > 3 {
				partIndex = 3
			}
			u.Browser.Version = part[partIndex].Version
			if engine.Name == "AppleWebKit" {
				switch part[partLen-1].Name {
				case "Edge":
				case "EdgA":
					u.Browser.Name = "Edge"
					u.Browser.Version = part[partLen-1].Version

					break
				case "Edg":
					u.Browser.Name = "Edge"
					u.Browser.Version = part[partLen-1].Version

					break
				case "OPR":
					u.Browser.Name = "Opera"
					u.Browser.Version = part[partLen-1].Version

					break
				default:
					switch part[partLen-2].Name {
					case "Electron":
						u.Browser.Name = "Electron"
						u.Browser.Version = part[partLen-2].Version
					case "DuckDuckGo":
						u.Browser.Name = "DuckDuckGo"
						u.Browser.Version = part[partLen-2].Version
					case "PhantomJS":
						u.Browser.Name = "PhantomJS"
						u.Browser.Version = part[partLen-2].Version
					default:
						switch part[partIndex].Name {
						case "Chrome", "CriOS":
							u.Browser.Name = "Chrome"
							break
						case "HeadlessChrome":
							u.Browser.Name = "Headless Chrome"
							break
						case "Chromium":
							u.Browser.Name = "Chromium"
							break
						case "GSA":
							u.Browser.Name = "Google App"
							break
						case "FxiOS":
							u.Browser.Name = "Firefox"
							break
						default:
							u.Browser.Name = "Safari"
						}
						u.Browser.Version = part[partIndex].Version
					}
				}
			} else if engine.Name == "Presto" {
				u.Browser.Name = "Opera"
				u.Browser.Version = u.Browser.IdVersion
			} else if engine.Name == "Gecko" {
				name := part[2].Name
				if name == "MRA" && partLen > 4 {
					name = part[4].Name
					u.Browser.Version = part[4].Version
				}
				u.Browser.Name = name
			} else if engine.Name == "like" && part[2].Name == "Gecko" {
				//兼容IE11老古董
				u.Browser.Engine = "Trident"
				u.Browser.Name = "Internet Explorer"
				for _, c := range part[0].Other {
					version := regexp.MustCompile("^rv:(.+)$").FindStringSubmatch(c)
					if len(version) > 0 {
						u.Browser.Version = version[1]
						return
					}
				}
				u.Browser.Version = ""
			}
		}
	} else if partLen == 1 && len(part[0].Other) > 1 {
		//兼容IE8、9、10老古董
		other := part[0].Other
		if other[0] == "compatible" && strings.HasPrefix(other[1], "MSIE") {
			u.Browser.Engine = "Trident"
			u.Browser.Name = "Internet Explorer"
			for _, v := range other {
				if strings.HasPrefix(v, "Trident/") {
					switch v[8:] {
					case "4.0":
						u.Browser.Version = "8.0"
						break
					case "5.0":
						u.Browser.Version = "9.0"
						break
					case "6.0":
						u.Browser.Version = "10.0"
						break
					}
					break
				}
			}
			//返回MSIE的版本号
			if u.Browser.Version == "" {
				u.Browser.Version = strings.TrimSpace(other[1][4:])
			}
		}
	}
}
