package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/astaxie/beego/httplib"
)

type MovieHTML struct {
	URL     string
	Content string
}

func GetMovieHTML(url string) MovieHTML {
	res := httplib.Get(url)
	str, err := res.String()
	if err != nil {
		panic(err)
	}
	if str == "" {
		panic("The Website is Empty")
	}
	Result := MovieHTML{}
	Result.URL = url
	Result.Content = str
	return Result
}

func (m *MovieHTML) GetUrls() []string {
	//<a href="https://movie.douban.com/subject/3578981/?from=subject-page" >
	urls := regexp.MustCompile(`<h3><a class="mdbColor f15" href="/mdb/film/(.*)/" target="_blank">.*</a></h3>`)
	result := urls.FindAllStringSubmatch(m.Content, -1)
	before := "https://www.1905.com/mdb/film/"
	var movielist []string
	for _, v := range result {
		movielist = append(movielist, before+v[1]+"/")
	}
	return movielist
	//return result[0]
}

func (m *MovieHTML) GetName() []string {
	//<a href="https://movie.douban.com/subject/3578981/?from=subject-page" >
	urls := regexp.MustCompile(`<h1>(.*).*?<span>.*</span>.*?<span class="score">`)
	result := urls.FindAllStringSubmatch(m.Content, -1)
	var movielist []string
	for _, v := range result {
		ase := strings.Split(v[1], " ")
		movielist = append(movielist, ase[0])
	}
	return movielist
	//return result[0]
}
func (m *MovieHTML) GetSocre() []string {
	//<a href="https://movie.douban.com/subject/3578981/?from=subject-page" >
	urls := regexp.MustCompile(`<span class="score"><b>(.*)</span>`)
	result := urls.FindAllStringSubmatch(m.Content, -1)
	var movielist []string
	for _, v := range result {
		di1 := strings.Split(v[1], "</b>")
		movielist = append(movielist, di1[0]+di1[1])
	}
	return movielist
	//return result[0]
}
func removeDuplicateElement(languages []string) []string {
	result := make([]string, 0, len(languages))
	temp := map[string]struct{}{}
	for _, item := range languages {
		if _, ok := temp[item]; !ok { //如果字典中找不到元素，ok=false，!ok为true，就往切片中append元素。
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func (m *MovieHTML) GetMType() []string {
	//<a href="http://www.1905.com/mdb/film/list/.*/?fr=mdbypsy_jsy_lx" target="_blank" data-hrefexp="fr=mdbypsy_jsy_lx">剧情</a>
	urls := regexp.MustCompile(`<a href=".*?" target="_blank" data-hrefexp="fr=mdbypsy_jsy_lx">(.*?)</a>`)
	result := urls.FindAllStringSubmatch(m.Content, -1)
	var movielist []string

	for _, v := range result {
		//di1 := strings.Split(v[1], "</b>")
		movielist = append(movielist, v[1])
	}
	if len(movielist) == 0 {
		return nil
	}
	movielist = removeDuplicateElement(movielist)
	fmt.Println(movielist)
	for v, k := range movielist {
		if strings.Contains(k, "&nbsp;") {
			fmt.Println(v, k)
			movielist = movielist[:v]
			break
		}
	}
	return movielist
	//return result[0]
}
func (m *MovieHTML) GetRegion() []string {
	//<a href="http://www.1905.com/mdb/film/list/country-China/" target="_blank" data-hrefexp="fr=mdbypsy_jsy_gj">中国</a>
	urls := regexp.MustCompile(`<a href="http://www.1905.com/mdb/film/list/country-.*/" target="_blank" data-hrefexp="fr=mdbypsy_jsy_gj">(.*?)</a>`)
	result := urls.FindAllStringSubmatch(m.Content, -1)
	if len(result) == 0 {
		return nil
	}
	k := result[0][0]
	reg2 := regexp.MustCompile(`[\p{Han}]+`)
	movielist := reg2.FindAllString(k, -1)
	return movielist
	//return result[0]
}

func (m *MovieHTML) GetDurationAndDate() []string {
	//<a href="https://movie.douban.com/subject/3578981/?from=subject-page" >
	urls := regexp.MustCompile(`</span>|<span class="information-item">(.*)</span>|`)
	result := urls.FindAllStringSubmatch(m.Content, -1)
	var movielist []string
	for _, v := range result {
		//di1 := strings.Split(v[1], "</b>")
		movielist = append(movielist, v[1])
	}
	movielist = removeDuplicateElement(movielist)
	movielist = movielist[1:]
	movielist[0] = strings.Split(movielist[0], "(内地)")[0]

	//fmt.Println(movielist[0])
	return movielist
	//return result[0]
}

func (m *MovieHTML) GetDir() []string {
	//<div class="creator-name">周显扬</div><span class="creator-class">导演</span>
	urls := regexp.MustCompile(`<a class="creator-resume .*" href="http://www.1905.com/mdb/star/.*?/" target="_blank" title="(.*)" data-hrefexp="fr=mdbypsy_zy">`)
	result := urls.FindAllStringSubmatch(m.Content, -1)
	var movielist []string
	for _, v := range result {
		//di1 := strings.Split(v[1], "</b>")
		movielist = append(movielist, v[1])
	}
	movielist = removeDuplicateElement(movielist)
	return movielist
	//return result[0]
}
func (m *MovieHTML) GetActor() []string {
	//<a data-hrefexp="fr=mdbypsy_zy" class="creator-resume" href=".*?" target="_blank" title="(.*)">
	urls := regexp.MustCompile(`<a data-hrefexp="fr=mdbypsy_zy" class="creator-resume" href=".*?" target="_blank" title="(.*)">`)
	result := urls.FindAllStringSubmatch(m.Content, -1)
	var movielist []string
	for _, v := range result {
		//di1 := strings.Split(v[1], "</b>")
		movielist = append(movielist, v[1])
	}
	movielist = removeDuplicateElement(movielist)
	if len(movielist) > 3 {
		movielist = movielist[:3]
	}
	return movielist
	//return result[0]
}

type I905 struct {
	Name   string
	Dir    []string
	Actors []string
	Time   string
	Dur    string
	Region []string
	mtype  []string
	Score  string
}

//...
func GetAll(url string) (I905, error) {
	var tar I905

	res := GetMovieHTML(url)
	//fmt.Println("GetMovieHTML Success")
	Name := res.GetName()
	if len(Name) == 0 {
		return tar, fmt.Errorf("WRONG")
	}
	tar.Name = Name[0]
	tar.Dir = res.GetDir()
	tar.Actors = res.GetActor()
	TandD := res.GetDurationAndDate()
	if len(TandD) != 2 {
		return tar, fmt.Errorf("WRONG")
	}
	tar.Time = TandD[0]
	tar.Dur = TandD[1]

	Region := res.GetRegion()
	if len(Region) == 0 {
		return tar, fmt.Errorf("WRONG")
	}
	tar.Region = Region
	tar.mtype = res.GetMType()
	Score := res.GetSocre()
	if len(Score) == 0 {
		return tar, fmt.Errorf("WRONG")
	}
	tar.Score = Score[0]
	return tar, nil
}
func URLGenerator() []string {
	var seeds []string
	//page := 1
	for i := 2019; i <= 2019; i++ {
		fmt.Println(fmt.Sprint(i))
		page := 10
		for page < 13 {
			seed := "https://www.1905.com/mdb/film/calendaryear/" + fmt.Sprint(i) + "#" + fmt.Sprint(page)
			seeds = append(seeds, seed)
			page++
		}

	}

	return seeds
}
func main1() {

	//movieurl := res.GetUrls()
	f, err := os.Create("haha3.csv")
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
	//url := URLGenerator()
	candidateUrl := URLGenerator()
	var URLS []string
	//url := "https://www.1905.com/mdb/film/calendaryear/2018#9"
	for _, k := range candidateUrl {
		res := GetMovieHTML(k)
		for _, m := range res.GetUrls() {
			URLS = append(URLS, m)
		}
	}
	//url := "https://www.1905.com/mdb/film/calendaryear/2019#6"
	for _, v := range URLS {
		res, err := GetAll(v)
		if err != nil {
			continue
		}

		fmt.Println(res)
		//w.Write([]string{res.Name, "", "", res.Time, res.Dur, "", res.Score})
		dirs := strings.Join(res.Dir, " ")
		act := strings.Join(res.Actors, " ")
		reg := strings.Join(res.Region, " ")
		mt := strings.Join(res.mtype, " ")
		w.Write([]string{res.Name, "", dirs, act, res.Time, res.Dur, reg, mt, res.Score})
		w.Flush()
	}

	fmt.Println("Finish")
}

func main() {
	main1()

}
