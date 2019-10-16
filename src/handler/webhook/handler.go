// @File     : handler
// @Author   : Ville
// @Time     : 19-10-16 上午9:44
// webhook
package webhook

import (
	"github.com/tidwall/gjson"
	"mo-gateway/src/handler/webhook/travis"
	"net/url"
	"strings"
	"vilgo/vlog"
)

type Handler struct {
}

func (sel *Handler) TravisPost(body []byte, args ...interface{}) (interface{}, error) {
	val, err := url.ParseQuery(string(body))
	if err != nil {
		vlog.ERROR("解析参数错误：%s", err.Error())
		return nil, err
	}
	if len(val) < 0 {
		vlog.ERROR("参数不存在")
		return "fail", nil
	}
	payload := val.Get("payload")
	if strings.Trim(payload, " ") == "" {
		return "fail", nil
	}
	result := gjson.Parse(payload)
	repository := result.Get("repository")
	title := ""
	if repository.Raw != "" {
		title = repository.Get("name").String()
	}
	return "success", travis.TraivsWeChat(
		title,
		result.Get("state").String(),
		result.Get("committer_email").String(),
		result.Get("type").String(),
	)
}
func (sel *Handler) TravisGet(map[string][]string, ...interface{}) (interface{}, error) {
	return nil, nil
}
