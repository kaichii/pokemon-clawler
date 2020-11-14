package main

import (
	"fmt"

	"github.com/antchfx/htmlquery"
	"github.com/pokemon-clawler/utils"
)

const baseUrl string = "https://wiki.52poke.com"

func main() {

	urls := make(chan string)

	go generateUrls(fmt.Sprintf("%s/wiki/%s", baseUrl, string("宝可梦列表（按全国图鉴编号）/简单版")), urls)

	parse(urls)
}

func generateUrls(url string, out chan<- string) {
	doc, err := htmlquery.LoadURL(url)

	utils.CheckError(err, "[load url]:")

	list, err := htmlquery.QueryAll(doc, "//table[@class=\"a-c roundy eplist bgl-神奇宝贝百科 b-神奇宝贝百科 bw-2\"]/tbody/tr/td[last()]/a[@class=\"mw-redirect\"]")

	utils.CheckError(err, "[query //tr/td[last()]/a[@class=\"mw-redirect\"]]:")

	for _, a := range list {
		out <- fmt.Sprintf("%s%s", baseUrl, htmlquery.SelectAttr(a, "href"))
	}

	close(out)
}

func parse(in <-chan string) {
	for url := range in {
		fmt.Println(url)
		// htmlquery.LoadURL(url)
	}
}
