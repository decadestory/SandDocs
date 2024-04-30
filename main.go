package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/decadestory/goutil/conf"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
)

type DocInfo struct {
	Date  string
	Title string
	Desc  string
	Path  string
}

func main() {
	cateHtml := GetCateHtml()
	GenCates(cateHtml)
	GenIndex(conf.Configs.GetString("index_page"), cateHtml)
}

func GenCates(cateHtml string) {
	filepath.WalkDir("docs/", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && path != "docs/" {
			GenCate(filepath.Base(path), cateHtml)
		}
		return nil
	})
}

func GenCate(path string, cateHtml string) {
	fmt.Println("start_genCate")

	tempData, _ := os.ReadFile("resource/template-doc.html")

	listHmml := GenDocs(path, cateHtml)
	cateHtml = strings.Replace(cateHtml, "../../docs_list/", "", -1)
	template := strings.Replace(string(tempData), "../../resource/", "../resource/", -1)
	template = strings.Replace(template, "{title}", path, -1)
	template = strings.Replace(template, "{cates}", cateHtml, -1)
	template = strings.Replace(template, "{content}", listHmml, -1)

	file, _ := os.Create("docs_list/" + path + ".html")
	defer file.Close()
	file.WriteString(template)

	fmt.Println("end_genCate")
}

func GenDocs(cateName, cateHtml string) string {

	var listHtml strings.Builder
	reg := regexp.MustCompile(`\[(.*)\]`)
	docList := []DocInfo{}

	filepath.Walk("docs/"+cateName+"/", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		name := filepath.Base(path)
		name = strings.TrimSuffix(name, ".md")

		title, desc := GenDoc(cateName+"/"+name, cateHtml)

		date := reg.FindString(desc)
		if date == "" {
			date = "[2000-01-01]"
		}
		date = strings.TrimPrefix(date, "[")
		date = strings.TrimSuffix(date, "]")
		desc = reg.ReplaceAllString(desc, "")
		docList = append(docList, DocInfo{Date: date, Title: title, Desc: desc, Path: path})

		return nil
	})

	slice.SortBy(docList, func(a, b DocInfo) bool {
		return a.Date > b.Date
	})

	for _, dv := range docList {
		li := fmt.Sprintf(`<div class="list-item">
		<div class="list-title"><a href="%s">%s</a></div>
		<div class="list-desc">%s</div>
		<div class="list-date">%s</div>
		</div>`, "../"+strings.Replace(dv.Path, ".md", ".html", 1), dv.Title, dv.Desc, dv.Date)
		listHtml.WriteString(li)
	}

	return listHtml.String()
}

func GenDoc(path string, cateHtml string) (string, string) {
	fmt.Println("start")

	tempData, _ := os.ReadFile("resource/template-doc.html")
	mdData, _ := os.ReadFile("docs/" + path + ".md")

	title := strings.SplitN(path, "/", 2)[1]
	desc := conf.Configs.GetString("no_desc")
	re := regexp.MustCompile(`<!--(.*)-->`)
	matchs := re.FindAllString(string(mdData), 1)
	if len(matchs) > 0 {
		desc, _ = strings.CutPrefix(matchs[0], "<!--")
		desc, _ = strings.CutSuffix(desc, "-->")
		desc = strings.TrimSpace(desc)
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(extension.Table),
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
				highlighting.WithFormatOptions(
					chromahtml.WithCustomCSS(map[chroma.TokenType]string{
						chroma.Line:    "padding-left: 10px;display:block;font-family:'Consolas','DejaVu Sans Mono','Bitstream Vera Sans Mono', monospace",
						chroma.Keyword: "font-style: normal;",
					}),
				),
			),
		),
	)

	var buf bytes.Buffer
	markdown.Convert(mdData, &buf)

	html := buf.String()
	template := strings.Replace(string(tempData), "{title}", path, -1)
	template = strings.Replace(template, "{cates}", cateHtml, -1)
	template = strings.Replace(template, "{content}", html, -1)

	file, _ := os.Create("docs/" + path + ".html")
	defer file.Close()
	file.WriteString(template)

	fmt.Println("end")
	return title, desc
}

func GenIndex(path string, cateHtml string) {
	fmt.Println("start")

	tempData, _ := os.ReadFile("resource/template-doc.html")
	mdData, _ := os.ReadFile("docs/" + path + ".md")

	markdown := goldmark.New(
		goldmark.WithExtensions(extension.Table),
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
				highlighting.WithFormatOptions(
					chromahtml.WithAllClasses(true),
				),
			),
		),
	)

	var buf bytes.Buffer
	markdown.Convert(mdData, &buf)

	html := buf.String()

	cateHtml = strings.Replace(cateHtml, "../../docs_list/", "docs_list/", -1)
	template := strings.Replace(string(tempData), "../../resource/", "resource/", -1)

	template = strings.Replace(template, "{title}", path, -1)
	template = strings.Replace(template, "{cates}", cateHtml, -1)
	template = strings.Replace(template, "{content}", html, -1)

	file, _ := os.Create("index.html")
	defer file.Close()
	file.WriteString(template)

	fmt.Println("end")
}

func GetCateHtml() string {
	cates := []string{}
	filepath.WalkDir("docs/", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && path != "docs/" {
			cates = append(cates, filepath.Base(path))
		}
		return nil
	})

	sort.Strings(cates)
	var sb strings.Builder
	for _, v := range cates {
		idx := strings.IndexByte(v, '.')
		if idx < 0 {
			idx = 0
		} else {
			idx++
		}
		sb.WriteString(fmt.Sprintf(`<li><a href="%s">%s</a></li>`, "../../docs_list/"+v+".html", v[idx:]))
	}

	return sb.String()
}
