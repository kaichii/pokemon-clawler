# pokemon clawler

## 前言

最近想弄点数据来练手，遂就想到从网上爬点儿 pokemon 的信息。刚好最近学 go，于是就想用 go 来实现，随便推荐 [go 语言圣经](https://books.studygolang.com/gopl-zh/) 这本书。这是我第一个爬 🐛 项目，实现的也很简单。源 [wiki.52poke.com](https://wiki.52poke.com/wiki/%E5%AE%9D%E5%8F%AF%E6%A2%A6%E5%88%97%E8%A1%A8%EF%BC%88%E6%8C%89%E5%85%A8%E5%9B%BD%E5%9B%BE%E9%89%B4%E7%BC%96%E5%8F%B7%EF%BC%89/%E7%AE%80%E5%8D%95%E7%89%88)， 解析 html 用到 [htmlquery](https://pkg.go.dev/github.com/antchfx/htmlquery@v1.2.3), [htmlquery](https://pkg.go.dev/github.com/antchfx/htmlquery@v1.2.3) 使用的是[xpath 选择器](https://www.w3school.com.cn/xpath/xpath_syntax.asp)。

github: [项目源码](https://github.com/kaichii/pokemon-clawler)

## [htmlquery](https://pkg.go.dev/github.com/antchfx/htmlquery@v1.2.3) 基本用法

### 常用方法

#### LoadURL

```go
// 根据给定的 url 返回该 url 的 HTML document
func LoadURL(url string) (*html.Node, error)
```

#### Find

```go
// 返回 top 节点下所有满足 `expr` 的所有节点
func Find(top *html.Node, expr string) []*html.Node
```

#### FindOne

```go
// 返回 top 节点下满足 `expr` 的第一个节点
func FindOne(top *html.Node, expr string) *html.Node
```

#### InnerText

```go
// 返回节点 tag 间的文本
func InnerText(n *html.Node) string
```

#### SelectAttr

```go
// 返回 n 节点上属性名为 name 的属性值
func SelectAttr(n *html.Node, name string) (val string)
```

## 爬取数据

### 生成要爬取的网页链接

根据给定的主页链接，获取需要获取的 pokemon 的主页链接，然后丢到 channel 里，等待后续的处理。

```go
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

```

### 提取数据

提取链接的 HTML 文本里我们想要获取的数据， 主要就是分析 HTML 文档结构，写出定位到我们需要其数据的文档节点的 xpath 表达式，提取数据就完事了。

处理 channel in 里链接， 并从中提取数据丢到 channel out 里。

```go
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
```

### 导出数据

```go
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
```