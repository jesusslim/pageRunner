package page_runner

import (
	"fmt"
	"github.com/jesusslim/slimmysql"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PageRunner struct {
	id        string
	waitGroup sync.WaitGroup
	sql       *slimmysql.Sql
	path      string    //扫描路径
	suffix    string    //文件后缀
	baseUrl   string    //http基础url
	ignore    []string  //忽略的目录/文件名
	rep       []string  //需替换为空的字符串
	extError  []string  //其他错误标示
	c         chan bool //channel 建议根据服务器情况调整合适大小
	result    map[string]map[string]interface{}
	cookie    string //cookie for sessionid
	maxTimes  int    //最大尝试访问次数
}

func NewPageRunner(id string, path string, suffix string, baseUrl string, ignore []string, rep []string, extErr []string, sql *slimmysql.Sql, channelLength int, cookie string, maxTimes int) *PageRunner {
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl = baseUrl + "/"
	}
	if suffix == "" {
		suffix = "index.html"
	}
	if channelLength == 0 {
		channelLength = 1
	}
	if maxTimes == 0 {
		maxTimes = 20
	}
	return &PageRunner{
		sql:      sql,
		path:     path,
		suffix:   suffix,
		baseUrl:  baseUrl,
		ignore:   ignore,
		extError: extErr,
		rep:      rep,
		c:        make(chan bool, channelLength),
		result:   make(map[string]map[string]interface{}),
		cookie:   cookie,
		maxTimes: maxTimes,
	}
}

//for thinkphp
func NewPageRunnerTP(id string, path, baseUrl string, rp []string, extErr []string, sql *slimmysql.Sql, channelLength int, cookie string, maxTimes int) *PageRunner {
	return NewPageRunner(
		id,
		path,
		".html",
		baseUrl,
		[]string{
			"Widget",
		},
		append(rp, "View/index.html", "View/"),
		extErr,
		sql,
		channelLength,
		cookie,
		maxTimes)
}

type UrlModel struct {
	module     string
	controller string
	action     string
	url        string
	times      int
	lastErr    string
}

func NewUrlModel(baseUrl, subUrl string) *UrlModel {
	mdl := &UrlModel{
		url:   baseUrl + subUrl,
		times: 0,
	}
	subs := strings.Split(subUrl, "/")
	l := len(subs)
	if l > 0 {
		mdl.module = subs[0]
	}
	if l > 1 {
		mdl.controller = subs[1]
	}
	if l > 2 {
		mdl.action = subs[2]
	}
	return mdl
}

func (this *PageRunner) fetchUrl(id string, url *UrlModel) {
	var start time.Time
	var err error
	var resp *http.Response
	if this.cookie == "" {
		//1.easy
		start = time.Now()
		resp, err = http.Get(url.url)
	} else {
		//2.with cookie
		client := &http.Client{}
		req, _ := http.NewRequest("GET", url.url, nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Cookie", this.cookie)
		start = time.Now()
		resp, err = client.Do(req)
	}
	duration := time.Since(start).Seconds() * 1000
	if err != nil {
		<-this.c
		fmt.Println("ERROR:" + err.Error())
		if url.times == this.maxTimes-1 {
			//结束 give up
			this.result[id] = map[string]interface{}{
				"duration":    0,
				"url":         url.url,
				"status":      -1,
				"create_time": time.Now().Unix(),
				"same":        0,
				"task_id":     this.id,
				"is_err":      1,
				"err":         url.lastErr,
				"module":      url.module,
				"controller":  url.controller,
				"action":      url.action,
				"times":       url.times + 1,
			}
			this.waitGroup.Done()
		} else {
			url.times++
			url.lastErr = err.Error()
			this.c <- true
			this.fetchUrl(id, url)
		}
	} else {
		defer resp.Body.Close()

		_, ok := this.result[id]

		if !ok {
			status := resp.StatusCode
			create_time := time.Now().Unix()
			header := resp.Header
			last_url := resp.Request.URL.String()
			same := 0
			is_err := 0
			if status != 200 {
				is_err = 1
			}
			if strings.EqualFold(url.url, last_url) {
				same = 1
			}
			// fmt.Println(header.Get("Date"))
			server := header.Get("Server")
			//xby := header.Get("X-Powered-By")
			if same == 0 {
				for _, e := range this.extError {
					if strings.Contains(last_url, e) {
						is_err = 1
						break
					}
				}
			}

			this.result[id] = map[string]interface{}{
				"duration":    int(duration),
				"url":         url.url,
				"status":      status,
				"create_time": create_time,
				"server":      server,
				"last_url":    last_url,
				"same":        same,
				"task_id":     this.id,
				"is_err":      is_err,
				"module":      url.module,
				"controller":  url.controller,
				"action":      url.action,
				"err":         url.lastErr,
				"times":       url.times + 1,
			}
		}

		<-this.c
		this.waitGroup.Done()
	}
}

func (this *PageRunner) walkDir() (map[string]*UrlModel, error) {
	files := map[string]*UrlModel{}
	suffix := strings.ToUpper(this.suffix)
	id := 0
	err := filepath.Walk(this.path, func(filename string, info os.FileInfo, err_inside error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(info.Name()), suffix) {
			for _, ig := range this.ignore {
				if strings.Contains(filename, ig) {
					return nil
				}
			}
			for _, rp := range this.rep {
				filename = strings.Replace(filename, rp, "", -1)
			}
			if strings.HasPrefix(filename, "/") {
				filename = filename[1:]
			}
			files[strconv.Itoa(id)] = NewUrlModel(this.baseUrl, filename)
			id++
		}
		return nil
	})
	return files, err
}

func (this *PageRunner) Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	urls, err := this.walkDir()
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	} else {
		num := len(urls)
		id, _ := this.sql.Table("task").Add(map[string]interface{}{
			"num":         num,
			"title":       "Test",
			"create_time": time.Now().Unix(),
		})
		this.id = string(strconv.FormatInt(id, 10))
		for k, v := range urls {
			this.waitGroup.Add(1)
			this.c <- true
			go this.fetchUrl(k, v)
		}
		this.waitGroup.Wait()
		fmt.Println("DataOk")
		success := len(this.result)
		this.sql.Clear().Table("task").Where("id = "+this.id).SetInc("success", success)
		this.sql = this.sql.Clear()
		for _, v := range this.result {
			this.sql.Table("url").Add(v)
		}
		fmt.Println("Finished")

	}
}

//example
// func main() {
// 	slimmysql.RegisterConnectionDefault(false, "127.0.0.1", "3307", "test", "root", "root", "", false)
// 	sql, _ := slimmysql.NewSqlInstanceDefault()
// 	baseUrl := "http://localhost:8888/teenager/Student"
// 	path := "/Applications/MAMP/htdocs/teenager/Application/Student"
// 	extErr := []string{
// 		"404",
// 		"error",
// 	}
// 	runner := NewPageRunnerTP("0", path, baseUrl,[]string{"/Applications/MAMP/htdocs/teenager/Application/"}, extErr, sql, 30, "PHPSESSID=c3146dcc95ba4e5992441718296aef1d", 50)
// 	runner.Run()
// }
