package main

import (
	"os"
	"bufio"
	"fmt"
	"strings"
	"regexp"
	"mdparser/src/htmltemplate"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

const footer = `</main>

<footer>
<hr>
<a href="https://github.com/Pippadi/MysoreLUGWebsite">Contribute on GitHub</a><br>
<a href="/">Read More</a>
</footer>
</body>
</html>`

func replaceCharAt(str, toInsert string, index int) string {
	return str[:index] + toInsert + str[index+1:]
}

func main() {
	var html = ""
	var template = htmltemplate.NewHTMLTemplate("templates")

	const SpecialCharacters = "_`()[]*"
	var SpecialCharacterNames = [7]string{"underaotuscore", "bac988utick", "openo88uhphesis", "closeoen3parenis", "opeqb38f5racket", "clo9342sqbrac", "ast8898erisk"}

	file, err := os.Open("template-article.md")
	checkError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	inParagraph := false
	currentParagraph := ""
	currentCodeBlk := ""
	inCode := false
	for scanner.Scan() {
		line := scanner.Text()

		for i, chr := range SpecialCharacters {
			line = strings.ReplaceAll(line, "\\"+string(chr), SpecialCharacterNames[i])
		}

		m := regexp.MustCompile(`\((.+?)\)\[(.+?)\]`) // Links as (text)[URL]
		line = m.ReplaceAllString(line, `<a href="${2}">${1}</a>`)

		m = regexp.MustCompile("`([^`]+?)`") // Inline code in ``
		line = m.ReplaceAllString(line, `<code>$1</code>`)

		m = regexp.MustCompile(`_(.+?)_`) // Italics in _..._
		line = m.ReplaceAllString(line, `<em>$1</em>`)

		m = regexp.MustCompile(`\*(.+?)\*`) // Bold in *...*
		line = m.ReplaceAllString(line, `<b>$1</b>`)

		if !inCode && line == "" {
			if inParagraph && currentParagraph != "" {
				template.AddParagraph(currentParagraph)
				currentParagraph = "" 
			}
			inParagraph = !inParagraph
		} else if line[:3] == "```" {
			if inCode {
				template.AddCodeBlk(currentCodeBlk[:len(currentCodeBlk) - 1])
				currentCodeBlk = ""
			}
			inCode = !inCode
		} else if line[:3] == "---" {
			template.AddHorizontalLine()
		} else if line[:2] == "# " {
			template.SetTitle(line[2:])
		} else if line[:3] == "## " {
			template.SetTitle(line[3:])
		} else {
			if inCode {
				currentCodeBlk += line + "\n"
			} else {
				currentParagraph += line + "\n"
			}
		}
	}

	template.Finalize()
	html += template.String()

	for i, chr := range SpecialCharacters {
			html = strings.ReplaceAll(html, SpecialCharacterNames[i], string(chr))
	}

	fmt.Println(html)
	checkError(scanner.Err())
}
