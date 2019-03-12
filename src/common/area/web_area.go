/***********************************************************************
* @ 获取地区信息，web接口
* @ brief
	、web api 不稳定，可能超时卡顿；查询、实测合适api
	、某些地区call不通，比如天朝的墙

* @ author zhoumf
* @ date 2019-3-12
***********************************************************************/
package area

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	kTaobaoLoacl = "http://ip.taobao.com/service/getIpInfo.php?ip=myip"
	kTaobaoApi   = "http://ip.taobao.com/service/getIpInfo.php?ip="
	kCountryApi  = "http://api.wipmania.com" //国内很卡
)

type TArea struct {
	Ip        string `json:"ip"`
	Country   string `json:"country"`
	Area      string `json:"area"`
	Region    string `json:"region"`
	CountryId string `json:"country_id"`
	AreaId    string `json:"area_id"`
	RegionId  string `json:"region_id"`
	Isp       string `json:"isp"`
	IspId     string `json:"isp_id"`
}
type retTaobao struct {
	Code int   `json:"code"`
	Data TArea `json:"data"`
}

func GetArea() (ret TArea) {
	if resp, err := http.Get(kTaobaoLoacl); err == nil {
		buf, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err == nil {
			var ack retTaobao
			json.Unmarshal(buf, &ack)
			ret = ack.Data
		}
	}
	return
}
func GetAreaEx(ip string) (ret TArea) {
	if resp, err := http.Get(kTaobaoApi + ip); err == nil {
		buf, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err == nil {
			var ack retTaobao
			json.Unmarshal(buf, &ack)
			ret = ack.Data
		}
	}
	return
}
func GetCountryId() string {
	if resp, err := http.Get(kCountryApi); err == nil {
		buf, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err == nil {
			s := string(buf)
			fmt.Println(s)
			return s
		}
	}
	return ""
}
