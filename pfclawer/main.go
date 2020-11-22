package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	htt "github.com/qy/httptest"
)

//MovieHTML 基本结构
type MovieHTML struct {
	URL     string
	Content string
}

//GetMovieHTML 获取网页源码
func GetMovieHTML(url string) MovieHTML {
	res := htt.GetWithCookies(url)

	if res == "" {
		panic("网页为空")
	}
	Result := MovieHTML{}
	Result.URL = url
	Result.Content = res
	return Result
}

//<span>页次：<span class="pager_number">1/15</span></span>
//GetIds 获取电影ID
func (m *MovieHTML) GetPages() []string {
	//<a href="/boxoffice/history/9291 6904" title="流浪地球">流浪地球</a>
	urls := regexp.MustCompile(`<span>页次：<span class="pager_number">1/(.*)</span></span>`)
	result := urls.FindAllStringSubmatch(m.Content, -1)
	//fmt.Println(result)
	var movielist []string
	for _, v := range result {
		movielist = append(movielist, v[1])
	}
	return movielist
	//return result[0]
}
func (m *MovieHTML) GetIds() []string {
	//<a href="/boxoffice/history/9291 6904" title="流浪地球">流浪地球</a>
	urls := regexp.MustCompile(`<a href="/boxoffice/history/(.*)" title=".*">.*</a>`)
	result := urls.FindAllStringSubmatch(m.Content, -1)
	//fmt.Println(result)
	var movielist []string
	for _, v := range result {
		movielist = append(movielist, v[1])
	}
	return movielist
	//return result[0]
}

//GetPF 获取电影票房
func (m *MovieHTML) GetPF() []string {
	//<h3 class="panel-title">阿涅斯·瓦尔达在这里和那里票房统计(最新票房 1.09万)</h3>
	urls := regexp.MustCompile(`<h3 class="panel-title">.*最新票房 (.*).*</h3>`)
	result := urls.FindAllStringSubmatch(m.Content, -1)
	//fmt.Println(result)
	var movielist []string
	for _, v := range result {
		movielist = append(movielist, v[1])
	}
	return movielist
	//return result[0]
}

//GetName 获取电影票房
func (m *MovieHTML) GetName() []string {
	//<h3 class="panel-title">阿涅斯·瓦尔达在这里和那里票房统计(最新票房 1.09万)</h3>
	urls := regexp.MustCompile(`<h3 class="panel-title">(.*)票房统计.*</h3>`)
	result := urls.FindAllStringSubmatch(m.Content, -1)
	//fmt.Println(result)
	var movielist []string
	for _, v := range result {
		movielist = append(movielist, v[1])
	}
	return movielist
	//return result[0]
}

//URLGenerator 总共27页
func URLGenerator() []string {
	var firstUrl []string
	for i := 2015; i <= 2019; i++ {
		candi := "http://58921.com/alltime/" + fmt.Sprint(i) + "?page="
		firstpage := "http://58921.com/alltime/" + fmt.Sprint(i)
		firstUrl = append(firstUrl, firstpage)
		prepage := GetMovieHTML(firstpage)
		pages := prepage.GetPages()[0]
		k, _ := strconv.Atoi(pages)
		p := 1
		for p <= k {
			firstUrl = append(firstUrl, candi+fmt.Sprint(p))
			p++
		}
		//fmt.Println(firstUrl)
	}

	return firstUrl
}

//URLbuilder 生成获取票房的页面地址
func URLbuilder(idlist []string) []string {
	candi1 := "http://58921.com/film/"
	candi2 := "/boxoffice"
	cand3 := new([]string)
	for _, v := range idlist {
		*cand3 = append(*cand3, candi1+v+candi2)
	}
	return *cand3
}

//PF 票房
func PF(url string) (string, float64) {
	target := GetMovieHTML(url)
	res := target.GetPF()
	names := target.GetName()
	if len(names) == 0 {
		return "", 0
	}
	val := strings.Split(res[0], ")")
	values := []rune(val[0])
	switch string(values[len(values)-1:]) {
	case "亿":
		floats, _ := strconv.ParseFloat(string(values[:len(values)-1]), 64)
		num := 1000000000 * floats
		return names[0], num
	case "万":
		floats, _ := strconv.ParseFloat(string(values[:len(values)-1]), 64)
		num := 10000 * floats
		return names[0], num
	default:
		floats, _ := strconv.ParseFloat(string(values), 64)
		return names[0], floats
	}
}

//MovieInformation 电影信息
type MovieInformation struct {
	Name  string
	Value float64
}

//All ...
func All() {
	f, err := os.Create("haha2.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f)
	w.Write([]string{
		"片名",
		"总票房",
		"导演",
		"主演",
		"上映时间",
		"片长",
		"制作国家",
		"类型",
		"评分(豆瓣电影)"})
	URLlist := URLGenerator()
	//URLlist = URLlist[2:]
	for _, v := range URLlist {
		//time.Sleep(5 * time.Second)
		fmt.Println(v)
		WebContent := GetMovieHTML(v)
		fmt.Println("Work on link" + v)
		MovieID := WebContent.GetIds()
		MovieContent := URLbuilder(MovieID)
		for _, k := range MovieContent {
			//time.Sleep(5 * time.Second)
			fmt.Println("Samll page: Sleep for 2 seconds")
			name, val := PF(k)
			if name == "" {
				continue
			}
			fmt.Println(name, val)
			w.Write([]string{name, fmt.Sprintln(int(val))})
			w.Flush()
		}
		fmt.Println("Finish " + v)
	}

	// res := GetMovieHTML(url)
	// list := res.GetIds()
	// res1 := URLbuilder(list)
	// fmt.Print(res1)
	//fmt.Println(PF(url))
}

func main() {
	//res := URLGenerator()
	//v := "http://58921.com/alltime/2015"
	//WebContent := GetMovieHTML(v)
	//fmt.Println(res)
	All()
}
