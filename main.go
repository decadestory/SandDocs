package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
)

func main() {
	cateHtml := GetCateHtml()
	GenDoc("HISTORY/Windows下visual studio code搭建golang开发环境", cateHtml)
	GenCates(cateHtml)
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
	cateHtml = strings.Replace(cateHtml, "../../docs_cate/", "", -1)
	template := strings.Replace(string(tempData), "../../resource/", "../resource/", -1)
	template = strings.Replace(template, "{title}", path, -1)
	template = strings.Replace(template, "{cates}", cateHtml, -1)
	template = strings.Replace(template, "{content}", listHmml, -1)

	file, _ := os.Create("docs_cate/" + path + ".html")
	defer file.Close()
	file.WriteString(template)

	fmt.Println("end_genCate")
}

func GenDocs(cateName, cateHtml string) string {

	var listHtml strings.Builder
	filepath.Walk("docs/"+cateName+"/", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		name := filepath.Base(path)
		idx := strings.IndexByte(name, '.')

		title, desc := GenDoc(cateName+"/"+name[:idx], cateHtml)
		li := fmt.Sprintf(`<div class="list-item">
		<div class="list-title">
		<a href="%s">%s</a> 
		</div><div class="list-desc">%s</div>
		</div>`, "../"+strings.Replace(path, ".md", ".html", 1), title, desc)
		listHtml.WriteString(li)
		return nil
	})

	return listHtml.String()
}

func GenDoc(path string, cateHtml string) (string, string) {
	fmt.Println("start")

	tempData, _ := os.ReadFile("resource/template-doc.html")
	mdData, _ := os.ReadFile("docs/" + path + ".md")

	title := strings.SplitN(path, "/", 2)[1]
	desc := "暂时没有描述"
	re := regexp.MustCompile(`<!-- (.*) -->`)
	matchs := re.FindAllString(string(mdData), 1)
	if len(matchs) > 0 {
		desc = strings.SplitN(matchs[0], " ", 2)[1]
		desc, _ = strings.CutPrefix(desc, "<!-- ")
		desc, _ = strings.CutSuffix(desc, " -->")
	}

	// Custom configuration
	markdown := goldmark.New(
		goldmark.WithExtensions(extension.Table),
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("paraiso-light"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
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

func GetCateHtml() string {
	cates := []string{}
	filepath.WalkDir("docs/", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && path != "docs/" {
			cates = append(cates, filepath.Base(path))
		}
		return nil
	})

	var sb strings.Builder
	for _, v := range cates {
		sb.WriteString(fmt.Sprintf(`<li><a href="%s">%s</a></li>`, "../../docs_cate/"+v+".html", v))
	}

	return sb.String()
}
