package htmltemplate

import (
	"path/filepath"
	"io/ioutil"
	"strings"
	"fmt"
)

type HTMLTemplate struct {
	header string
	footer string
	subtitle string
	para string
	codeBlk string
	horizontalLine string
	heading1 string
	heading2 string
	image string
	imageCaption string

	html string
	title string
	subtitles []string
	paragraphs []string
}

func NewHTMLTemplate(path string) *HTMLTemplate {
	h := HTMLTemplate{title:"PageTitle"}
	h.ReadTemplates(path)
	return &h
}

func (h *HTMLTemplate) ReadTemplates(path string) {
	content, err := ioutil.ReadFile(filepath.Join(path, "header"))
	h.header = string(content)
	h.html = h.header
	content, err = ioutil.ReadFile(filepath.Join(path, "footer"))
	h.footer = string(content)
	content, err = ioutil.ReadFile(filepath.Join(path, "subtitle"))
	h.subtitle= string(content)
	content, err = ioutil.ReadFile(filepath.Join(path, "para"))
	h.para= string(content)
	content, err = ioutil.ReadFile(filepath.Join(path, "codeBlk"))
	h.codeBlk= string(content)
	content, err = ioutil.ReadFile(filepath.Join(path, "horizontalLine"))
	h.horizontalLine = string(content)
	content, err = ioutil.ReadFile(filepath.Join(path, "heading1"))
	h.heading1 = string(content)
	content, err = ioutil.ReadFile(filepath.Join(path, "heading2"))
	h.heading2 = string(content)
	content, err = ioutil.ReadFile(filepath.Join(path, "image"))
	h.image = string(content)
	content, err = ioutil.ReadFile(filepath.Join(path, "imageCaption"))
	h.imageCaption = string(content)
	checkError(err)
}

func (h *HTMLTemplate) SetTitle(aTitle string) {
	h.html = strings.Replace(h.html, h.title, aTitle, 2)
	h.title = aTitle
}

func (h *HTMLTemplate) AddSubtitle(aSubtitle string) {
	h.html += strings.Replace(h.subtitle, "{}", aSubtitle, 1)
}

func (h *HTMLTemplate) AddHeading1(aHeading string) {
	h.html += strings.Replace(h.heading1, "{}", aHeading, 1)
}

func (h *HTMLTemplate) AddHeading2(aHeading string) {
	h.html += strings.Replace(h.heading2, "{}", aHeading, 1)
}

func (h *HTMLTemplate) AddParagraph(aParagraph string) {
	h.html += strings.Replace(h.para, "{}", aParagraph, 1)
}

func (h *HTMLTemplate) AddCodeBlk(aCodeBlk string) {
	h.html += strings.Replace(h.codeBlk, "{}", aCodeBlk, 1)
}

func (h *HTMLTemplate) AddHorizontalLine() {
	h.html += h.horizontalLine
}

func (h *HTMLTemplate) AddImage(path, opts string, captions []string) {
	caps := ""
	for _, cap := range captions {
		caps += strings.Replace(h.imageCaption, "{}", cap, 1)
	}
	blk := h.image
	blk = strings.Replace(blk, "{path}", path, 1)
	blk = strings.Replace(blk, "{opts}", opts, 1)
	blk = strings.Replace(blk, "{captions}", caps, 1)
	h.html += blk
}

func (h *HTMLTemplate) Finalize() {
	h.html += h.footer
}

func (h *HTMLTemplate) String() string {
	return h.html
}

func checkError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
