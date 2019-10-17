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
        "content": "### {{.Title}}  
State：<font color={{.color}}>{{.State}}</font>
Name：<font color=blue>{{.Name}}</font>
Email：<font color=blue>{{.Email}}</font>
ID：<font color=comment>{{.ID}}</font>
Number：<font color=warning>{{.Number}}</font>
Type：<font color=warning>{{.Type}}</font>
Start：<font color=comment>{{.StartTime}}</font>
End：<font color=comment>{{.EndTime}}</font>
Branch：<font color=comment>{{.Branch}}</font>
Message：<font color=comment>{{.Message}}</font>
Commit:<font color=comment>{{.Commit}}</font>
        "
    }
}
`

// http://hbimg.b0.upaiyun.com/828e72e2855b51a005732f4e007c58c92417a61e1bcb1-61VzNc_fw658
var weChatNoticTemplatePic = `
{
    "msgtype": "news",
    "news": {
       "articles" : [
			{
               "title" : "{{.Title}}",
               "description" : "{{.Message}}",
               "url" : "URL",
               "picurl" : "http://b-ssl.duitang.com/uploads/item/201808/27/20180827043223_twunu.jpg"
           	}
        ]
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
	// fmt.Println(string(msgBuf.String()))
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
