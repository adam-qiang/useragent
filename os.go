/**
 * Created by goland.
 * User: adam_wang
 * Date: 2023-07-14 00:36:02
 */

package useragent

import "strings"

// Os
// @Description: 系统结构体
type Os struct {
	Platform string
	Name     string
	Version  string
}

// 解析系统信息
// @receiver u *UserAgent
// @param p Part
func (u *UserAgent) analysisOs(p Part) {
	olen := len(p.Other)

	u.Os.Platform = getPlatformName(p.Other)
	if p.Name == "Mozilla" {
		if u.Os.Platform == "Windows" {
			u.Os.Name = getWindowsOsName(p.Other[0])
		} else {
			//其他操作系统需要配合浏览器引擎进行判断
			switch u.Browser.Engine {
			case "Gecko":
				//使用Firefox
				getGeckoOsName(u, p.Other)
			case "AppleWebKit":
				//使用AppleWebKit
				getWebkitOsName(u, p.Other)
			case "Trident":
				//使用IE
				getTridentOsName(u, p.Other)
			}
		}
	} else if p.Name == "Dalvik" {
		//安卓设备使用Dalvik虚拟机运行浏览器发送请求
		getDalvikOsName(u, p.Other)
	} else if p.Name == "Opera" {
		//Opera浏览器一个非开源浏览器（现在应该很少有人在用了）
		if olen > 0 {
			if strings.HasPrefix(p.Other[0], "Android") {
				u.IsMobile = true
			}
			u.Os.Name = p.Other[0]
		}
	}

	//根据系统名称解析版本
	osName := u.Os.Name
	if osName != "" {
		os := strings.Replace(osName, "like Mac OS X", "", 1)
		os = strings.Replace(os, "CPU", "", 1)
		os = strings.Trim(os, " ")

		//将系统名称进行切片操作
		osSplit := strings.Split(os, " ")

		//windows xp 64位版本的特殊处理
		if os == "Windows XP x64 Edition" {
			osSplit = osSplit[:len(osSplit)-2]
		}

		//将上面得到的切片系统名称进行拆分得到版本机器纯名称（即去除版本号的名称）
		u.Os.Name, u.Os.Version = splitOsName(osSplit)
	}

	//如果是移动设备则获取设备名称
	if u.IsMobile {
		u.getMobileInfo(p)
	}
}

// 获取平台名称
// @param Other []string
// @return string
func getPlatformName(Other []string) string {
	if len(Other) > 0 {
		if Other[0] != "compatible" {
			if strings.HasPrefix(Other[0], "Windows") {
				return "Windows"
			} else if strings.HasPrefix(Other[0], "Linux") || strings.HasPrefix(Other[0], "Android") {
				return "Linux"
			}

			return Other[0]
		}
	}
	return ""
}

// 获取windows平台的系统名称
// @param name string
// @return string
func getWindowsOsName(name string) string {
	sp := strings.SplitN(name, " ", 3)
	if len(sp) != 3 || sp[1] != "NT" {
		return name
	}

	switch sp[2] {
	case "5.0":
		return "Windows 2000"
	case "5.01":
		return "Windows 2000, Service Pack 1 (SP1)"
	case "5.1":
		return "Windows XP"
	case "5.2":
		return "Windows XP x64 Edition"
	case "6.0":
		return "Windows Vista"
	case "6.1":
		return "Windows 7"
	case "6.2":
		return "Windows 8"
	case "6.3":
		return "Windows 8.1"
	case "10.0":
		return "Windows 10"
	}
	return name
}

// Firefox浏览器解析系统名称
// @param p *UserAgent
// @param Other []string
func getGeckoOsName(u *UserAgent, Other []string) {
	olen := len(Other)
	if olen > 1 {
		if Other[1] == "U" || Other[1] == "arm_64" {
			if len(Other) > 2 {
				u.IsMobile = true
				u.Os.Name = Other[2]
			} else {
				u.Os.Name = Other[1]
			}
		} else {
			if strings.Contains(u.Os.Name, "Android") {
				u.IsMobile = true
				u.Os.Name = Other[1]
			} else if Other[0] == "Mobile" || Other[0] == "Tablet" {
				u.IsMobile = true
				u.Os.Name = "FirefoxOS"
			} else {
				u.Os.Name = Other[1]
			}
		}
	}
}

// Safari浏览器解析系统名称
// @param p *UserAgent
// @param Other []string
func getWebkitOsName(p *UserAgent, Other []string) {
	if p.Os.Platform == "Linux" {
		p.IsMobile = true
		if p.Browser.Name == "Safari" {
			p.Browser.Name = "Android"
		}
		if len(Other) > 1 {
			if Other[1] == "U" || Other[1] == "arm_64" {
				if len(Other) > 2 {
					p.Os.Name = Other[2]
				} else {
					p.IsMobile = false
					p.Os.Name = Other[0]
				}
			} else if len(Other) > 2 && Other[2] == "HarmonyOS" {
				//鸿蒙系统
				p.Os.Name = Other[2]
			} else {
				p.Os.Name = Other[1]
			}
		}
	} else if strings.HasPrefix(Other[0], "Windows NT") {
		p.Os.Name = getWindowsOsName(Other[0])
	} else if p.Os.Platform == "Macintosh" && p.Browser.Engine == "AppleWebKit" && p.Browser.Name == "Firefox" {
		//Ipad，在Firefox上成为Macintosh
		p.Os.Name = "iPad"
		p.IsMobile = true
	} else {
		p.Os.Name = Other[1]
	}
}

// IE浏览器解析系统名称（IE只运行在windows上）
// @param u *UserAgent
// @param Other []string
func getTridentOsName(u *UserAgent, Other []string) {
	olen := len(Other)

	if u.Os.Name == "" {
		u.Os.Name = "Windows"
		if olen > 2 {
			u.Os.Name = getWindowsOsName(Other[2])
		} else {
			u.Os.Name = "Windows NT 4.0"
		}
	}

	//windows phone设备
	for _, v := range Other {
		if strings.HasPrefix(v, "IEMobile") {
			u.IsMobile = true
			return
		}
	}
}

// 安卓设备使用Dalvik虚拟机进行访问解析系统名称
// @param p *UserAgent
// @param Other []string
func getDalvikOsName(u *UserAgent, Other []string) {
	olen := len(Other)

	if strings.HasPrefix(Other[0], "Linux") {
		u.Os.Name = Other[0]
		if olen > 2 {
			u.Os.Name = Other[2]
		}
		u.IsMobile = true
	}
}

// 拆分系统名称切片
// @param osSplit []string
// @return name string
// @return version string
func splitOsName(osSplit []string) (name string, version string) {
	if len(osSplit) == 1 {
		name = osSplit[0]
		version = ""
	} else {
		//一般情况下版本在系统名称后面
		nameSplit := osSplit[:len(osSplit)-1]
		version = osSplit[len(osSplit)-1]

		//Max OS X
		if len(nameSplit) >= 2 && nameSplit[0] == "Intel" && nameSplit[1] == "Mac" {
			nameSplit = nameSplit[1:]
		}
		name = strings.Join(nameSplit, " ")

		if strings.Contains(version, "x86") || strings.Contains(version, "i686") {
			//x86、i686是CPU架构不是版本所以得忽略掉
			version = ""
		} else if version == "X" && name == "Mac OS" {
			//Mac OS
			name = name + " " + version
			version = ""
		}
	}
	return name, version
}

// 获取移动设备信息
// @receiver u *UserAgent
// @param p Part
func (u *UserAgent) getMobileInfo(p Part) {
	olen := len(p.Other)
	if !u.IsMobile {
		return
	}
	if u.Os.Platform == "iPhone" || u.Os.Platform == "iPad" {
		u.Mobile = u.Os.Platform
		return
	}

	//Linux平台设备
	if p.Name == "Mozilla" && u.Os.Platform == "Linux" && olen > 2 {
		tmpSplit := make([]string, 0)
		for key, other := range p.Other {
			if strings.Contains(other, u.Os.Name) {
				//设备名称后面可能会追加地区标识需要忽略掉
				if strings.Contains(p.Other[key+1], "-") {
					tmpSplit = strings.Split(p.Other[key+1], "-")
				} else if strings.Contains(p.Other[key+1], "/") {
					tmpSplit = strings.Split(p.Other[key+1], "/")
				} else if strings.Contains(p.Other[key+1], "_") {
					tmpSplit = strings.Split(p.Other[key+1], "_")
				}
				if len(tmpSplit) == 2 && len(tmpSplit[0]) == 2 && len(tmpSplit[1]) == 2 {
					continue
				}

				u.Mobile = p.Other[key+1]

				break
			} else if strings.Contains(other, "Build") {
				u.Mobile = p.Other[key]
				break
			}
		}

		if u.Mobile != "" {
			tmp := strings.Split(u.Mobile, "Build")
			u.Mobile = strings.Trim(tmp[0], " ")

			return
		}
	}

	for _, v := range p.Other {
		if strings.Contains(v, "Build") {
			tmp := strings.Split(v, "Build")
			u.Mobile = strings.Trim(tmp[0], " ")

			return
		}
	}
}
