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
ID：<font color=comment>{{.ID}}</font>
Number：<font color=comment>{{.Number}}</font>
Name：<font color=comment>{{.Name}}</font>
Email：<font color=comment>{{.Email}}</font>
Type：<font color=comment>{{.Type}}</font>
State：<font color={{.color}}>{{.State}}</font>
Start：<font color=comment>{{.StartTime}}</font>
End：<font color=comment>{{.EndTime}}</font>
Branch：<font color=comment>{{.Branch}}</font>
Message：<font color=comment>{{.Message}}</font>
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

func TravisWeChat(params map[string]string) (err error) {
	tmpl := template.New("")
	tmpl.Parse(weChatHookTeamplate)
	state := params["State"]
	color := "green"
	if strings.ToLower(state) != "passed" {
		color = "red"
	}
	params["color"] = color
	msgBuf := bytes.NewBufferString("")
	// 执行模板字段替换
	if err = tmpl.Execute(msgBuf, params); err != nil {
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
