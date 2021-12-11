package main

import (
	"os"
	"bufio"
	"fmt"
	"strings"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

const header = `<html>
<head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta name="format-detection" content="address=no, telephone=no, date=no">
	<meta name="description" content="A website for Linux users and geeks in general to express their thoughts and spread information">
	<title>PageTitle</title>
	<link rel="stylesheet" type="text/css" href="../shared/css/main.css">
</head>
<body>
<header>
<h2 class="main-title">Mysore Linux Users' Group</h2>
<p class="main-subtitle">A website for Linux users and geeks in general to express their thoughts and spread information</p>
<hr>
</header>

<main>
`

const footer = `</main>

<footer>
<hr>
<a href="https://github.com/Pippadi/MysoreLUGWebsite">Contribute on GitHub</a><br>
<a href="/">Read More</a>
</footer>
</body>
</html>`

type articleHeading struct {
	Title     string
	Subtitles []string
}

func (h articleHeading) HTMLString() string {
	var str string
	str += `<h1 class="article-title">` + h.Title + "</h1>\n"
	for _, sub := range h.Subtitles {
		str += `<p class="article-subtitle">` + sub + "</p>\n"
	}
	return str
}

func main() {
	var heading = articleHeading{}
	var html = header

	file, err := os.Open("template-article.md")
	checkError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	headingNotDone := true
	for headingNotDone && scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, " ")

		if words[0] == "#" {
			heading.Title = line[2:]
		} else if words[0] == "##" {
			heading.Subtitles = append(heading.Subtitles, line[3:])
		} else if words[0] == "---" {
			headingNotDone = false
			html += heading.HTMLString()
			html += "<hr>\n"
		}
	}
	html = strings.Replace(html, "PageTitle", heading.Title, 1)

	inParagraph := false
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, " ")

		if line == "" && inParagraph {
			html += "</p>\n"
			inParagraph = false
		} else if words[0] == "###" {
			html += "<h3 class=\"article-heading1\">" + line[4:] + "</h3>\n"
		} else if words[0] == "---" {
			html += "<hr>\n"
		} else if line != "" {
			if !inParagraph {
				html += "\n<p class=\"article-paragraph\">\n"
				inParagraph = true
			}
			html += line + "\n"
		}
	}
	html += footer

	fmt.Println(html)
	checkError(scanner.Err())
}
