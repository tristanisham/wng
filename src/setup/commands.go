package setup

import (
	"bufio"
	"bytes"
	"encoding/json"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
)

func Init(path string) error {
	if path != "." {
		if err := os.Mkdir(path, 0775); err != nil {
			return err
		}
	}

	if err := os.Chdir("./" + path); err != nil {
		return err
	}
	// Program is now working on path. All future paths are relative to working directory.
	for _, dir := range []string{"src/assets", "src/posts"} {
		if err := os.MkdirAll(dir, 0775); err != nil {
			return err
		}
	}

	config_file := NewDefaultBlog()
	json, err := json.MarshalIndent(config_file, "", "    ")
	if err != nil {
		return err
	}

	if err := os.WriteFile("blog.json", json, 0775); err != nil {
		return err
	}

	if err := os.WriteFile("src/index.css", []byte(config_file.GenStyle()), 0775); err != nil {
		return err
	}

	if err := os.WriteFile("src/index.html", []byte(config_file.GenIndex()), 0775); err != nil {
		return err
	}

	return nil
}

func Build() (*DefaultBlog, error) {
	blog := new(DefaultBlog)
	cwd, _ := os.Getwd()
	os.Setenv("WRK_DIR", cwd)
	config, err := os.ReadFile("./blog.json")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(config, blog); err != nil {
		return nil, err
	}

	if err := blog.buildArticles(); err != nil {
		return nil, err
	}

	return blog, nil
}

func (b *DefaultBlog) buildArticles() error {
	b.Articles = make([]Article, 0)
	// if err := os.Chdir("src/posts"); err != nil {
	// 	return err
	// }
	if err := filepath.WalkDir("./src/posts", func(path string, d fs.DirEntry, err error) error {
		if strings.ContainsAny(d.Name(), "md") {
			article := new(Article)
			article.Slug = strings.Split(d.Name(), ".")[0]
			readFile, err := os.Open(path)
			if err != nil {
				return err
			}

			filescanner := bufio.NewScanner(readFile)
			filescanner.Split(bufio.ScanLines)
			var fileLines []string

			for filescanner.Scan() {
				fileLines = append(fileLines, filescanner.Text())
			}

			readFile.Close()

			b.parseOptions(fileLines, article)

			b.Articles = append(b.Articles, *article)
		}
		return nil
	}); err != nil {
		return err
	}

	return b.offload()
}

func (b *DefaultBlog) parseOptions(raw []string, article *Article) {
	options := make([]string, 0)
	body := make([]string, 0)
	for i, line := range raw {
		if line == "~~~" {
			options = append(options, raw[:i]...)
			body = append(body, raw[i+1:]...)
			break
		}
	}

	article.Body = string(blackfriday.Run([]byte(strings.Join(body, "\n"))))

	for _, line := range options {
		opts := strings.Split(line, ":")
		if len(opts) >= 2 {
			switch strings.ToLower(opts[0]) {
			case "title":
				article.Title = opts[1]
			case "subtitle":
				article.Subtitle = opts[1]
			case "public":
				choice := strings.TrimSpace(opts[1])
				if choice == "1" || choice == "true" {
					article.Public = true
				}
			case "tags", "keywords":
				article.Tags = strings.Split(strings.ReplaceAll(opts[1], "\"", ""), ",")
				for i := range article.Tags {
					article.Tags[i] = strings.TrimSpace(article.Tags[i])
					if len(article.Tags[i]) == 0 {
						article.Tags[i] = article.Tags[len(article.Tags)-1]
						article.Tags = article.Tags[:len(article.Tags)-1]
					}
				}
			}
		}

	}

}

func (b *DefaultBlog) offload() error {
	data, err := json.MarshalIndent(b, "", "    ")
	if err != nil {
		return err
	}

	root := os.Getenv("WRK_DIR")

	return os.WriteFile(root+"/blog.json", data, 0775)
}

func (b *DefaultBlog) Dist() error {
	for _, dir := range []string{"dist/assets"} {
		if err := os.MkdirAll(dir, 0775); err != nil {
			return err
		}
	}
	index, err := os.ReadFile("./src/index.html")
	if err != nil {
		return err
	}

	assets, err := os.ReadDir("src/assets")
	if err != nil {
		return err
	}


	for _, file := range assets {
		data, err := os.ReadFile(file.Name())
		if err != nil {
			return err
		}

		os.WriteFile("dist/assets/" + file.Name(), data, 0755)
	}

	css, err := os.ReadFile("src/index.css")
	if err != nil {
		return err
	}
	if err := os.WriteFile("dist/index.css", css, 0775); err != nil {
		return err
	}

	templ := template.Must(template.New("index").Parse(string(index)))
	

	for i := range b.Articles {
		b.Articles[i].HTML = template.HTML(b.Articles[i].Body)
	}

	var out bytes.Buffer
	templ.Lookup("index").Execute(&out, b)

	
	return os.WriteFile("dist/index.html", out.Bytes(), 0775)
	

}
