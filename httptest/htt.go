package htt

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func httpHandle(method, urlVal, data string) string {
	client := &http.Client{}
	var req *http.Request

	if data == "" {
		urlArr := strings.Split(urlVal, "?")
		if len(urlArr) == 2 {
			urlVal = urlArr[0] + "?" + getParseParam(urlArr[1])
		}
		req, _ = http.NewRequest(method, urlVal, nil)
	} else {
		req, _ = http.NewRequest(method, urlVal, strings.NewReader(data))
	}

	//添加cookie，key为X-Xsrftoken，value为df41ba54db5011e89861002324e63af81
	//可以添加多个cookie
	cookie1 := &http.Cookie{Name: "Hm_lvt_e71d0b417f75981e161a94970becbb1b", Value: "1605262002,1605280237,1605327029"}
	cookie2 := &http.Cookie{Name: "Hm_lpvt_e71d0b417f75981e161a94970becbb1b", Value: "1605327029"}
	cookie3 := &http.Cookie{Name: "DIDA642a4585eb3d6e32fdaa37b44468fb6c", Value: "1ota1sni7n76k4l41lnkj7jsh2"}
	cookie4 := &http.Cookie{Name: "remember", Value: "MTEzNTI2LjIxNjM0Mi4xMDI4MTYuMTA3MTAwLjExMTM4NC4yMDc3NzQuMTE5OTUyLjExMTM4NC4xMDI4MTYuMA%3D%3D"}
	cookie5 := &http.Cookie{Name: "time", Value: "MTEzNTI2LjIxNjM0Mi4xMDI4MTYuMTA3MTAwLjExMTM4NC4yMDc3NzQuMTE5OTUyLjExMTM4NC4xMDQ5NTguMTE1NjY4LjEwMjgxNi4xMTM1MjYuMTA5MjQyLjEwOTI0Mi4xMTM1MjYuMTA5MjQyLjEyMjA5NC4xMDI4MTYuMA%3D%3D"}
	req.AddCookie(cookie1)
	req.AddCookie(cookie2)
	req.AddCookie(cookie3)
	req.AddCookie(cookie4)
	req.AddCookie(cookie5)

	//添加header，key为X-Xsrftoken，value为b6d695bbdcd111e8b681002324e63af81
	//req.Header.Add("X-Xsrftoken", "b6d695bbdcd111e8b681002324e63af81")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	return string(b)
}

//将get请求的参数进行转义
func getParseParam(param string) string {
	return url.PathEscape(param)
}

//GetWithCookies ...
func GetWithCookies(url string) string {
	return httpHandle("GET", url, "")
}
