package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/pokemon-clawler/utils"
)

const baseUrl string = "https://wiki.52poke.com"

func main() {

	result, err := os.Create("pokemon.csv")

	utils.CheckError(err, "[creat file]:")

	defer result.Close()

	result.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(result)

	urls := make(chan string)
	columns := make(chan []string)

	go generateUrls(fmt.Sprintf("%s/wiki/%s", baseUrl, string("宝可梦列表（按全国图鉴编号）/简单版")), urls)

	go parse(urls, columns)

	for c := range columns {
		err := writer.Write(c)

		utils.CheckError(err, "[write line]:")
	}

	writer.Flush()
}

func generateUrls(url string, out chan<- string) {
	doc, err := htmlquery.LoadURL(url)

	utils.CheckError(err, "[load url]:")

	nodes, err := htmlquery.QueryAll(doc, "//table[@class=\"a-c roundy eplist bgl-神奇宝贝百科 b-神奇宝贝百科 bw-2\"]/tbody/tr/td[last()]/a[@class=\"mw-redirect\"]")

	utils.CheckError(err, "[query //tr/td[last()]/a[@class=\"mw-redirect\"]]:")

	for _, a := range nodes {
		out <- fmt.Sprintf("%s%s", baseUrl, htmlquery.SelectAttr(a, "href"))
	}

	close(out)
}

func parse(in <-chan string, out chan<- []string) {
	for url := range in {

		log.Println(url)

		name := strings.TrimLeft(url, baseUrl+"/wiki/")

		doc, err := htmlquery.LoadURL(url)

		utils.CheckError(err, "[load url]:")

		root := htmlquery.FindOne(doc, "//div[@class=\"mw-parser-output\"]/table[2]/tbody")

		result := []string{}

		order := ""

		chineseName := ""

		imageUri := ""

		description := ""

		orderNode := htmlquery.FindOne(root, "//a[@title=\"宝可梦列表（按全国图鉴编号）\"]")

		chineseNameNode := htmlquery.FindOne(root, "//td/span[@style=\"font-size:1.5em\"]/b")

		imageNode := htmlquery.FindOne(root, "/tr[2]//a[@class=\"image\"]/img")

		descriptionNode := htmlquery.FindOne(doc, "//div[@class=\"mw-parser-output\"]/p[2]")

		if chineseNameNode != nil {
			chineseName = htmlquery.InnerText(chineseNameNode)
		}

		if imageNode != nil {
			imageUri = "https:" + htmlquery.SelectAttr(imageNode, "data-url")
		}

		if descriptionNode != nil {
			description = strings.Replace(htmlquery.InnerText(descriptionNode), "\n", "", -1)
		}

		if orderNode != nil {
			order = htmlquery.InnerText(orderNode)
		}

		result = append(result, order, name, chineseName, imageUri, description)

		out <- result
	}

	close(out)
}
