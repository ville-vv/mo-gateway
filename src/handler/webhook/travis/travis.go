// @File     : travis
// @Author   : Ville
// @Time     : 19-10-16 下午3:20
// travis
package travis

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"vilgo/vlog"
)

var weChatHookTeamplate = `
{
    "msgtype": "markdown",
    "markdown": {
        "content": "#### {{.Title}}\n
提  交  人：<font color=comment>{{.Who}}</font>
事件类型：<font color=comment>{{.Event}}</font>
状    态：<font color={{.Color}}>{{.State}}</font>
        "
    }
}
`
var cli = http.Client{
	Transport: &http.Transport{
		//跳过证书验证
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func TraivsWeChat(title, state, who, event string) (err error) {
	tmpl := template.New("")
	tmpl.Parse(weChatHookTeamplate)

	color := "green"
	if strings.ToLower(state) != "passed" {
		color = "red"
	}
	msgBuf := bytes.NewBufferString("")
	// 执行模板字段替换
	if err = tmpl.Execute(msgBuf, map[string]string{"Title": title, "Who": who, "Event": event, "Color": color, "State": state}); err != nil {
		vlog.ERROR("创建模板失败：")
		return
	}

	reqHttp, err := http.NewRequest("POST", "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=02dd2384-04ae-426f-a85b-3b7884d0cfbe", msgBuf)
	if err != nil {
		vlog.ERROR("创建HTTP请求失败：%s", err.Error())
	}
	reqHttp.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(reqHttp)
	if err != nil {
		vlog.ERROR("请求微信机器人失败：%s", err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		vlog.ERROR("请求微信机器人失败 status：%s", resp.Status)
		return
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	vlog.INFO("请求微信机器人成功 status：%s", string(respBody))
	return nil
}
