/**
 * Created by goland.
 * User: adam_wang
 * Date: 2023-07-14 00:33:24
 */

package useragent

import (
	"strings"
)

// UserAgent
// @Description: User-Agent结构体
type UserAgent struct {
	Ua       string
	Os       Os
	Browser  Browser
	Mobile   string
	IsMobile bool
	IsRobot  bool
}

// Part
// @Description: User-Agent各部分信息结构体
type Part struct {
	Name    string
	Version string
	Other   []string
}

// New 创建一个新的User-Agent
// @param ua string
// @return *UserAgent
func New(ua string) *UserAgent {
	u := &UserAgent{}
	u.initialize()
	u.Ua = ua
	u.analysis(ua)

	return u
}

// 初始化User-Agent
// @receiver u *UserAgent
func (u *UserAgent) initialize() {
	u.Ua = ""
	u.Os.Platform = ""
	u.Os.Name = ""
	u.Os.Version = ""
	u.Browser.Id = ""
	u.Browser.IdVersion = ""
	u.Browser.Name = ""
	u.Browser.Version = ""
	u.Browser.Engine = ""
	u.Browser.EngineVersion = ""
	u.Mobile = ""
	u.IsMobile = false
	u.IsRobot = false

}

// analysis 解析User-Agent
// @receiver u *UserAgent
// @param ua string
func (u *UserAgent) analysis(ua string) {
	var part []Part
	for index, limit := 0, len(ua); index < limit; {
		analysis := analysisPart(ua, &index)
		if !u.IsMobile && analysis.Name == "Mobile" {
			u.IsMobile = true
		}
		part = append(part, analysis)
	}

	if len(part) > 0 {
		u.Browser.Id = part[0].Name
		u.Browser.IdVersion = part[0].Version

		u.analysisBrowser(part)
		//系统设备信息在浏览器标识符后面的括号中，即第一部分信息的Other中
		u.analysisOs(part[0])
	}
}

// 解析User-Agent各部分信息
// @param ua string
// @param index *int
// @return p Part
func analysisPart(ua string, index *int) (p Part) {
	// 根据User-Agent各部分信息放置顺序位置为依据，通过循环User-Agent整个信息对比不同信息放置规则进行匹配
	// 浏览器标识：通常位于User-Agent的开头的字符串，用于识别所使用的浏览器。在早期互联网的时候，Netscape Navigator是最流行的浏览器，其他浏览器为了兼容性通常会在其User-Agent字符串中包含“Mozilla”这个关键字，使得网站认为使用的是Netscape Navigator浏览器。后来随着Internet Explorer等其他浏览器的兴起这个惯例就一直被沿用下来了，现在很多浏览器在User-Agent字符串中仍然会包含这个关键字，以确保对旧的网站的兼容性
	// 操作系统信息：通常位于浏览器标识后面，一些操作系统信息可能带有附加的位数标识，如 "Win64" 表示64位的Windows系统，用于标识用户设备所运行的操作系统类型和版本。常见的操作系统信息包括Windows、Macintosh、Linux、Android、HarmonyOS（这个系统有点特殊）
	// 浏览器信息：位于操作系统信息后面，浏览器信息通常通过斜杠或空格与操作系统信息分隔，用于提供关于浏览器的具体信息：名称、版本号以及其他特定标识。常见的浏览器信息包括 Chrome、Safari、Firefox、Edge
	// 插件和功能信息：通常出现在User-Agent中的方括号（[]）中，用于指示特定的插件、功能或附加组件的信息，这些信息可能包括渲染引擎（如 "KHTML", "Gecko" 或 "WebKit"）、插件（如 "Adobe Acrobat"）和其他功能（如 "Mobile" 表示移动设备）

	var buffer []byte

	//检查是否是空的，符合条件后根据空格进行分割，且当前位置字符不是括号类的字符
	if *index < len(ua) && ua[*index] != '(' && ua[*index] != '[' {
		buffer = readUa(ua, index, ' ', false)
		p.Name, p.Version = partInfo(buffer)
	}

	//匹配其他信息
	if *index < len(ua) && ua[*index] == '(' {
		*index++
		buffer = readUa(ua, index, ')', true)
		p.Other = strings.Split(string(buffer), "; ")
		*index++
	}

	//[]里面的内容一般是插件和功能信息在这里我们不需要进行过滤
	if *index < len(ua) && ua[*index] == '[' {
		*index++
		_ = readUa(ua, index, ']', true)
		*index++
	}
	return p
}

// 读取User-Agent信息
// @param ua string
// @param index *int
// @param delimiter byte 指定分隔符
// @param isIgnoreNest bool 是否忽略嵌套在()、[]中的括号
// @return []byte
func readUa(ua string, index *int, delimiter byte, isIgnoreNest bool) []byte {
	var buffer []byte

	i := *index
	nestLan := 0
	for ; i < len(ua); i = i + 1 {
		if ua[i] == delimiter {
			if nestLan == 0 {
				*index = i + 1
				return buffer
			}
			nestLan--
		} else if isIgnoreNest && ua[i] == '(' {
			nestLan++
		}
		buffer = append(buffer, ua[i])
	}
	*index = i + 1
	return buffer
}

// 获取User-Agent各部分信息
// @param product []byte
// @return string
// @return string
func partInfo(product []byte) (string, string) {
	prod := strings.SplitN(string(product), "/", 2)
	if len(prod) == 2 {
		if len(prod) == 3 {
			return prod[1], prod[2]
		}
		return prod[0], prod[1]
	}
	return string(product), ""
}
