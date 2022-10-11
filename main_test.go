package main

import (
	"fmt"
	"testing"
)

func TestGetHttpHtmlContent(t *testing.T) {
	url := "https://xskydata.jobs.feishu.cn/school/?keywords=&category=&location=&project=&type=&job_hot_flag=&current=1&limit=10&functionCategory="

	content, err := GetHttpHtmlContent(url, ".title__bb7170.positionItem-title.sofiaBold")
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
}

func TestGetSpecialData(t *testing.T) {
	url := "https://xskydata.jobs.feishu.cn/school/?keywords=&category=&location=&project=&type=&job_hot_flag=&current=1&limit=10&functionCategory="

	content, err := GetHttpHtmlContent(url, title)
	if err != nil {
		panic(err)
	}
	data, err := GetSpecialData(content, title, city, typ, item, desc)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}
