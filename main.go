package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
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

type GenInfo struct {
	Changes         []string
	ChangeCates     []string
	ChangeFileNames []string
	CateHtml        string
}

func main() {
	gi := &GenInfo{}
	gi.ComputeChanges()
	gi.ComputeChangeCates()
	gi.ComputeChangeFiles()
	gi.GetCateHtml()
	gi.GenCates()
	gi.GenIndex()
}

func (gi *GenInfo) GenCates() {
	for _, change := range gi.ChangeCates {
		gi.GenCate(change)
	}
}

func (gi *GenInfo) GenCate(path string) {
	fmt.Println("start_genCate")

	tempData, _ := os.ReadFile("resource/template-doc.html")
	listHmml := gi.GenDocs(path)
	if listHmml == "" {
		return
	}

	cateHtml := strings.Replace(gi.CateHtml, "../../docs_list/", "", -1)
	template := strings.Replace(string(tempData), "../../resource/", "../resource/", -1)
	template = strings.Replace(template, "{title}", path, -1)
	template = strings.Replace(template, "{cates}", cateHtml, -1)
	template = strings.Replace(template, "{content}", listHmml, -1)

	file, _ := os.Create("docs_list/" + path + ".html")
	defer file.Close()
	file.WriteString(template)

	fmt.Println("end_genCate")
}

func (gi *GenInfo) GenDocs(cateName string) string {

	var listHtml strings.Builder
	reg := regexp.MustCompile(`\[(.*)\]`)
	docList := []DocInfo{}

	filepath.Walk("docs/"+cateName+"/", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		name := filepath.Base(path)
		name = strings.TrimSuffix(name, ".md")

		title, desc := gi.GenDoc(cateName + "/" + name)
		if title == "" || desc == "" {
			return nil
		}

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

func (gi *GenInfo) GenDoc(path string) (string, string) {
	fmt.Println("start")

	if !slices.Contains(gi.ChangeFileNames, filepath.Base(path+".md")) {
		return "", ""
	}

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
	template = strings.Replace(template, "{cates}", gi.CateHtml, -1)
	template = strings.Replace(template, "{content}", html, -1)

	file, _ := os.Create("docs/" + path + ".html")
	defer file.Close()
	file.WriteString(template)

	fmt.Println("end")
	return title, desc
}

func (gi *GenInfo) GenIndex() {
	fmt.Println("start")
	path := conf.Configs.GetString("index_page")
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

	cateHtml := strings.Replace(gi.CateHtml, "../../docs_list/", "docs_list/", -1)
	template := strings.Replace(string(tempData), "../../resource/", "resource/", -1)

	template = strings.Replace(template, "{title}", path, -1)
	template = strings.Replace(template, "{cates}", cateHtml, -1)
	template = strings.Replace(template, "{content}", html, -1)

	file, _ := os.Create("index.html")
	defer file.Close()
	file.WriteString(template)

	fmt.Println("end")
}

func (gi *GenInfo) GetCateHtml() {

	var sb strings.Builder
	for _, v := range gi.ChangeCates {
		idx := strings.IndexByte(v, '.')
		if idx < 0 {
			idx = 0
		} else {
			idx++
		}
		sb.WriteString(fmt.Sprintf(`<li><a href="%s">%s</a></li>`, "../../docs_list/"+v+".html", v[idx:]))
	}

	gi.CateHtml = sb.String()
}

func (gi *GenInfo) ComputeChanges() {

	oldFileMd5s := map[string]string{}
	newFileMd5s := map[string]string{}

	filepath.WalkDir("docs/", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && filepath.Ext(path) == ".md" {

			path = strings.ReplaceAll(path, "\\", "/")
			fileMd5, _ := FileMD5(path)
			newFileMd5s[path] = fileMd5
		}
		return nil
	})

	jsonPath := "conf/file_md5.json"
	if _, err := os.Stat(jsonPath); err == nil {
		jsonData, _ := os.ReadFile(jsonPath)
		json.Unmarshal(jsonData, &oldFileMd5s)
	}

	file, _ := os.Create(jsonPath)
	jsonData, _ := json.Marshal(newFileMd5s)
	file.Write(jsonData)
	file.Close()

	if len(oldFileMd5s) == 0 {
		all := []string{}
		for k := range newFileMd5s {
			all = append(all, k)
		}
		sort.Strings(all)
		gi.Changes = all
		return
	}

	res := []string{}
	for k, v := range newFileMd5s {
		if val, ok := oldFileMd5s[k]; ok && val != v {
			res = append(res, k)
		}
	}

	sort.Strings(res)
	gi.Changes = res

}

func (gi *GenInfo) ComputeChangeCates() {
	changeCates := []string{}
	for _, change := range gi.Changes {
		cate := strings.Split(change, "/")
		changeCates = append(changeCates, cate[1])
	}
	sort.Strings(changeCates)
	gi.ChangeCates = changeCates
}

func (gi *GenInfo) ComputeChangeFiles() {
	changeFileNames := []string{}
	for _, change := range gi.Changes {
		changeFileNames = append(changeFileNames, filepath.Base(change))
	}
	sort.Strings(changeFileNames)
	gi.ChangeFileNames = changeFileNames
}

func FileMD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
