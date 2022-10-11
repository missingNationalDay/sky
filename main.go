package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	json "github.com/json-iterator/go"
	"log"
	"os"
	"strings"
	"time"
)

var s = "https://xskydata.jobs.feishu.cn/school/?keywords=&category=&location=&project=&type=&job_hot_flag=&current=%d&limit=10&functionCategory="

const (
	//岗位名称
	title = ".title__bb7170.positionItem-title.sofiaBold"
	//城市
	city = ".subTitle__bb7170.positionItem-subTitle > span:nth-child(1)"
	typ  = ".subTitle__bb7170.positionItem-subTitle > span:nth-child(3) > span"
	item = ".subTitle__bb7170.positionItem-subTitle > span:nth-child(5) > span"
	//工作描述
	desc = ".jobDesc__bb7170.positionItem-jobDesc"
)

type Job struct {
	//岗位名称
	Title string
	//岗位城市
	City string
	//岗位类型
	Typ string
	//招聘项目
	Item string
	//介绍
	Desc string
}

func GetHttpHtmlContent(url string, selector string, sel interface{}) (string, error) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	//初始化参数，先传一个空的数据
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	// create context
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	// 执行一个空task, 用提前创建Chrome实例
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	//创建一个上下文，超时时间为40s
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 40*time.Second)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector),
		chromedp.OuterHTML(`document.querySelector("body")`, &htmlContent, chromedp.ByJSPath),
	)
	if err != nil {
		//logger.Info("Run err : %v\n", err)
		fmt.Println("run err", err)
		return "", err
	}
	//log.Println(htmlContent)

	return htmlContent, nil
}

func GetSpecialData(htmlContent string, title, city, typ, item, desc string) ([]Job, error) {
	jobs := make([]Job, 10)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	dom.Find(title).Each(func(i int, selection *goquery.Selection) {
		jobs[i].Title = selection.Text()
	})
	dom.Find(city).Each(func(i int, selection *goquery.Selection) {
		jobs[i].City = selection.Text()
	})
	dom.Find(typ).Each(func(i int, selection *goquery.Selection) {
		jobs[i].Typ = selection.Text()
	})
	dom.Find(item).Each(func(i int, selection *goquery.Selection) {
		jobs[i].Item = selection.Text()
	})
	dom.Find(desc).Each(func(i int, selection *goquery.Selection) {
		jobs[i].Desc = selection.Text()
	})
	//fmt.Println(jobs)
	return jobs, nil
}

func main() {
	datas := make([]Job, 0)
	for i := 1; i <= 5; i++ {
		url := fmt.Sprintf(s, 1)
		content, err := GetHttpHtmlContent(url, title, "")
		if err != nil {
			panic(err)
		}
		data, err := GetSpecialData(content, title, city, typ, item, desc)
		if err != nil {
			panic(err)
		}
		datas = append(datas, data...)
	}

	j, err := json.MarshalToString(datas)
	if err != nil {
		panic(err)
	}
	file, err := os.Create("jobs.json")
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(j)
	defer file.Close()
	if err != nil {
		panic(err)
	}

}
