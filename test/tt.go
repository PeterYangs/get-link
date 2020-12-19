package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"time"
)

var ws fyne.Window

func main() {

	a := app.New()
	ws = a.NewWindow("网站链接爬取")

	//w.SetFixedSize(true)

	//fyne.TextWrap(1)

	input := widget.NewEntry()
	input.SetPlaceHolder("请输入网址")

	input.Wrapping = fyne.TextWrapWord

	input.MultiLine = true

	label := widget.NewLabel("")

	//label.Wrapping=fyne.TextTruncate
	label.Wrapping = fyne.TextWrapBreak

	c := container.NewVBox(
		//hello,
		input,
		label,
		widget.NewButton("click", func() {

			label.SetText("sssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss")
		}),
	)

	//c.Resize(fyne.Size{Width:400})

	ws.SetContent(c)

	//设置窗体大小
	ws.Resize(fyne.Size{Width: 400, Height: 300})

	ws.ShowAndRun()

	//fmt.Println("gg")

}

func reSize() {

	for {

		ws.Resize(fyne.Size{Height: 300, Width: 400})

		time.Sleep(time.Microsecond * 500)
	}

}
