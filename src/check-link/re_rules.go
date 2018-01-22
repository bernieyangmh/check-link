package check_link

import "regexp"

const (
	PATTERN_SRC   = `src=\"(.*?)\"`
	PATTERN_HERF  = `href=\"(.*?)\"`
	PATTERN_HTTP  = `^http(.*?)`
	PATTERN_LINK  = `^https?:\/\/[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`
	PATTERN_SLASH = `^/(.*?)`
	ALLOW_DOMAIN  = `(qiniu.com)|(qiniu.com.cn)`
)


//返回匹配href=的相对路径数组
func ReDomainMatch(s string) bool {
	reDomain, _ := regexp.Compile(ALLOW_DOMAIN)
	isAllow := reDomain.MatchString(s)
	return isAllow

}

//返回匹配href=的相对路径数组
func ReHrefSubMatch(s string) [][]string {
	reHref, _ := regexp.Compile(PATTERN_HERF)
	hrefArray := reHref.FindAllStringSubmatch(s, 10000)

	return hrefArray

}

//返回匹配src=的相对路径数组
func ReSrcSubMatch(s string) [][]string {
	reHref, _ := regexp.Compile(PATTERN_SRC)
	srcArray := reHref.FindAllStringSubmatch(s, 10000)

	return srcArray

}

func ReLinkSubMatch(s string) [][]string {
	reLink, _ := regexp.Compile(PATTERN_LINK)
	srcArray := reLink.FindAllStringSubmatch(s, 10000)

	return srcArray

}

//re匹配http链接
func ReIsHttp(s string) bool {
	reHttp, _ := regexp.Compile(PATTERN_HTTP)
	a := reHttp.MatchString(s)
	return a
}

//re匹配链接
func ReIsLink(s string) bool {
	reLink, _ := regexp.Compile(PATTERN_LINK)
	a := reLink.MatchString(s)

	return a

}

//re匹配slasp
func ReHaveSlash(s string) bool {
	reSlash, _ := regexp.Compile(PATTERN_SLASH)
	a := reSlash.MatchString(s)

	return a

}
