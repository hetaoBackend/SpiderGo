package spider

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var urlFormat = "http://180.96.8.44/WindSearch/WindSearch/handler/SearchHandler.ashx?index=newsearchresult&version=1&pageIndex=1&pageNum=25&key=%v&dimension=2018.HK&mtype=5&type=-2&newsflg=0&btflg=-1&rtflg=0&userTypes=&range=0&sp=title&wind.sessionid=06a54262cd534a38a0197673f0752bf6&fileds=title,windcode,spell,fullspell,tagcode,areacode,keyword,content,abstract&sort=_score%20desc,publishdate%20desc&suggest=0&suggestfields=title,windcode,spell,fullspell,tagcode,areacode,section,keyword,author,content&d=212558&rppright=1"

func Spider(searchInfo string, wg *sync.WaitGroup, fileLock *sync.Mutex) { //生成client 参数为默认

	client := &http.Client{}
	//生成要访问的url
	searchInfo = strings.Replace(searchInfo, " ", "", -1)
	url := fmt.Sprintf(urlFormat, searchInfo)
	//提交请求
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		wg.Done()
		panic(err)
	}
	//增加header选项
	request.Header.Add("Cookie", "wind.sessionid=06a54262cd534a38a0197673f0752bf6; ASP.NET_SessionId=2k00ly553pdvukauqsvbyavx")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/605.1.15 (KHTML, like Gecko)")

	//处理返回结果
	response, _ := client.Do(request)
	fmt.Println(response.StatusCode)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body) //请求数据进行读取
	if err != nil {
		fmt.Println("err:%v", err)
		wg.Done()
		return
	}
	res, err := simplejson.NewJson(body)

	if err != nil {
		fmt.Println(searchInfo)
		fmt.Printf("%v\n", err)
		wg.Done()
		return
	}
	rows, _ := res.Get("dataor").String()
	finalRes, err := simplejson.NewJson([]byte(rows))
	if err != nil {
		fmt.Printf("%v\n", err)
		wg.Done()
		return
	}
	codeList, _ := finalRes.Get("hits").Get("hits").Array()
	for _, v := range codeList {
		if source, ok := v.(map[string]interface{})["_source"]; ok {
			if codeType, ok := v.(map[string]interface{})["_type"]; !ok || codeType != "bond" {
				continue
			}
			if tagCode, ok := source.(map[string]interface{})["tagcode"]; ok {
				fmt.Println(searchInfo, tagCode)
				fileLock.Lock()
				tracefile(fmt.Sprintf("%v %v", searchInfo, tagCode.([]interface{})[0]))
				fileLock.Unlock()
				if err != nil {
					fmt.Printf("%v\n", err)
				}
				break
			}
		}
	}
	wg.Done()
}

func tracefile(str_content string) {
	fd, _ := os.OpenFile("a.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fd_time := time.Now().Format("2006-01-02 15:04:05")
	fd_content := strings.Join([]string{"======", fd_time, "=====", str_content, "\n"}, "")
	buf := []byte(fd_content)
	fd.Write(buf)
	fd.Close()
}
