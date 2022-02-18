package main

import (
	"os"
	"bufio"
	"fmt"
	"mdparser/src/htmltemplate"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: ./mdparser markdown.md [templates folder]")
		os.Exit(1)
	}
	var templatePath = "templates"
	if len(os.Args) > 2 {
		templatePath = os.Args[2]
	}
	var template = htmltemplate.NewHTMLTemplate(templatePath)

	file, err := os.Open(os.Args[1])
	checkError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

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
			} else if line[:1] == "!" {
				path := line[1:]
				scanner.Scan()
				opts := ""
				l := scanner.Text()
				captions := make([]string, 0)
				for l != "" && scanner.Scan() {
					if l[0] == '!' {
						opts += l[1:] + " "
					} else {
						captions = append(captions, l)
					}
					l = scanner.Text()
				}
				template.AddImage(path, opts, captions)
			} else {
				scanner.Scan()
				para := line + "\n"
				l := scanner.Text()
				for l != "" && scanner.Scan() {
					para += l + "\n"
					l = scanner.Text()
				}
				template.AddParagraph(para)
			}
		}
	}

	template.Finalize()

	fmt.Print(template.String())
	checkError(scanner.Err())
}
