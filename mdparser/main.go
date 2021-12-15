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
	if len(os.Args) == 1 {
		fmt.Println("Usage: ./mdparser markdown.md [templates folder]")
		os.Exit(1)
	}
	var html = ""
	var templatePath = "templates"
	if len(os.Args) > 2 {
		templatePath = os.Args[2]
	}
	var template = htmltemplate.NewHTMLTemplate(templatePath)

	const SpecialCharacters = "_`()[]*"
	var SpecialCharacterNames = [7]string{"underaotuscore", "bac988utick", "openo88uhphesis", "closeoen3parenis", "opeqb38f5racket", "clo9342sqbrac", "ast8898erisk"}

	file, err := os.Open(os.Args[1])
	checkError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")

		for i, chr := range SpecialCharacters {
			line = strings.ReplaceAll(line, "\\"+string(chr), SpecialCharacterNames[i])
		}

		m := regexp.MustCompile(`\[(.+?)\]\((.+?)\)`) // Links as (text)[URL]
		line = m.ReplaceAllString(line, `<a href="${2}">${1}</a>`)

		m = regexp.MustCompile("`([^`]+?)`") // Inline code in ``
		line = m.ReplaceAllString(line, `<code>$1</code>`)

		m = regexp.MustCompile(`_(.+?)_`) // Italics in _..._
		line = m.ReplaceAllString(line, `<em>$1</em>`)

		m = regexp.MustCompile(`\*(.+?)\*`) // Bold in *...*
		line = m.ReplaceAllString(line, `<b>$1</b>`)

		if line != "" {
			if line[:3] == "```" {
				blk := ""
				scanner.Scan()
				l := scanner.Text()
				for l != "```" && scanner.Scan() {
					blk += l + "\n"
					l = scanner.Text()
				}
				template.AddCodeBlk(blk[:len(blk) - 1])
			} else if line[:3] == "---" {
				template.AddHorizontalLine()
			} else if line[:2] == "# " {
				template.SetTitle(line[2:])
			} else if line[:3] == "## " {
				template.AddSubtitle(line[3:])
			} else if line[:4] == "### " {
				template.AddHeading1(line[4:])
			} else if line[:5] == "#### " {
				template.AddHeading2(line[5:])
			} else if line[0] == '!' {
				path := line[1:]
				scanner.Scan()
				opts := ""
				l := strings.Trim(scanner.Text(), " \t")
				captions := make([]string, 0)
				for l != "" && scanner.Scan() {
					if l[0] == '!' {
						opts += l[1:] + " "
					} else {
						captions = append(captions, l)
					}
					l = strings.Trim(scanner.Text(), " \t")
				}
				template.AddImage(path, opts, captions)
			} else {
				scanner.Scan()
				para := line
				l := strings.Trim(scanner.Text(), " \t")
				for l != "" && scanner.Scan() {
					para += l + "\n"
					l = strings.Trim(scanner.Text(), " \t")
				}
				template.AddParagraph(para)
			}
		}
	}

	template.Finalize()
	html += template.String()

	for i, chr := range SpecialCharacters {
			html = strings.ReplaceAll(html, SpecialCharacterNames[i], string(chr))
	}

	fmt.Print(html)
	checkError(scanner.Err())
}
