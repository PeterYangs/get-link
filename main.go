package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/PeterYangs/tools"
	"github.com/PuerkitoBio/goquery"
	"github.com/flopp/go-findfont"
	"os"
	"regexp"
	"strings"
	"sync"
)

//存放链接
var list = make(map[string]string)

var mutex sync.Mutex

var wait sync.WaitGroup

//最大递归层数（防止内存溢出）
var maxLv = 50

//提交按钮
var submitButton *widget.Button

//停止按钮
var stopButton *widget.Button

//是否点击停止按钮
var isStop = false

//是否已完成
var isDone = false

//日志文本框
var log *widget.Label

//滚动框
var Scroll *container.Scroll

func main() {
	a := app.New()
	w := a.NewWindow("网站链接爬取")

	//设置窗体大小
	w.Resize(fyne.Size{Width: 400, Height: 200})

	input := widget.NewEntry()
	input.SetPlaceHolder("请输入网址")

	//提交按钮
	submitButton = widget.NewButton("提交", func() {

		go do(input.Text)

	})

	stopButton = widget.NewButton("停止", func() {

		//fmt.Println("dd")

		isStop = true

		for {

			if isDone {

				done()

				break

			}

		}

	})

	stopButton.Hide()

	//c:=

	log = widget.NewLabel("")

	//log=

	Scroll = container.NewVScroll(log)

	//Scroll.Resize(fyne.Size{Height: 400,Width:400})

	//设置滚动框的最新尺寸
	Scroll.SetMinSize(fyne.Size{Height: 400, Width: 400})

	//Scroll.

	w.SetContent(container.NewVBox(
		//hello,
		input,
		submitButton,
		stopButton,
		Scroll,
	))

	//c.A

	w.ShowAndRun()

	os.Unsetenv("FYNE_FONT")
}

func get() {

	re, _ := tools.GetWithString("https://www.925g.com")

	fmt.Println(re)

}

func init() {
	//获取中文字体列表
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		//设置字体
		if strings.Contains(path, "simkai.ttf") {
			os.Setenv("FYNE_FONT", path)
			break
		}
	}
}

func do(urls string) {

	//提交按钮禁用
	submitButton.Disable()

	//显示停止按钮
	stopButton.Show()
	//stopButton.E

	host := urls

	url := urls

	//删除文件
	tools.DeleteFile(getFileNameWithHost(host))

	for i := 0; i < 10; i++ {

		wait.Add(1)

		go getUrl(url, host, &wait, 1)

	}

	wait.Wait()

	fmt.Println("执行完毕")

	isDone = true

	done()

	//恢复按钮
	submitButton.Enable()

}

func getUrl(url string, host string, wait *sync.WaitGroup, lv int) {

	if url == "" {

		return
	}

	defer func() {

		if wait != nil {

			wait.Done()
		}

	}()

	if isStop {

		return

	}

	if lv >= maxLv {

		return
	}

	html, err := tools.GetWithString(url)

	if err != nil {

		//log.Fatal(err)

		fmt.Println(err)

		return

	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)

		return
	}

	doc.Find("a").Each(func(i int, selection *goquery.Selection) {

		h, _ := selection.Attr("href")

		//fmt.Println(h)

		if h == "/" {

			return

		}

		findHttp, _ := regexp.MatchString("^"+host, h)

		findLocal, _ := regexp.MatchString(`^/.*`, h)

		if findHttp {

			mutex.Lock()

			_, ok := list[h]

			//fmt.Println(ok)

			mutex.Unlock()

			if ok {

				//fmt.Println("存在")

			} else {

				//fmt.Println(h)

				//log.Add(widget.NewLabel(h))

				log.SetText(log.Text + h + "\n")

				Scroll.ScrollToBottom()

				//getUrl(h,host)

				mutex.Lock()

				list[h] = ""

				//panic(strings.Replace(host,"/","_",-1))

				tools.WriteLine(getFileNameWithHost(host), h)

				mutex.Unlock()

				getUrl(h, host, nil, lv+1)

			}

		}

		if findLocal {

			mutex.Lock()

			_, ok := list[host+h]

			mutex.Unlock()

			if ok {

			} else {

				//fmt.Println(host + h)

				//log.Add(widget.NewLabel(host + h))

				log.SetText(log.Text + host + h + "\n")

				Scroll.ScrollToBottom()

				mutex.Lock()

				list[host+h] = ""

				tools.WriteLine(getFileNameWithHost(host), host+h)

				mutex.Unlock()

				getUrl(host+h, host, nil, lv+1)

			}

		}

	})

}

func getFileNameWithHost(host string) string {

	return strings.Replace(strings.Replace(strings.Replace(host, "/", "_", -1), ":", "", -1), ".", "_", -1) + ".txt"

}

//任务完成
func done() {

	list = make(map[string]string)

	mutex = sync.Mutex{}

	submitButton.Enable()

	isDone = false

	log.SetText("")

	stopButton.Hide()

}
