package setup

import "html/template"

type DefaultBlog struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Keywords    []string  `json:"keywords"`
	Theme       string    `json:"theme"`
	Articles    []Article `json:"articles"`
}

func NewDefaultBlog() DefaultBlog {
	return DefaultBlog{
		Title:       "",
		Description: "",
		Keywords:    []string{},
		Theme:       "",
		Articles:     []Article{},
	}
}

type Templater interface {
	GenIndex() string
}

//GenStyle returns the default css file as a string.
func (d *DefaultBlog) GenStyle() string {
	return `@media screen and (min-width: 800px){body{background-color:#f0f0f0}article{padding-left:20%;padding-right:20%;padding-top:5%}}@media screen and (max-width: 799px){article{padding:2%}}@media screen{h1{font-size:30pt}p{font-size:16pt;line-height:2}code{font-size:16pt}li{font-size:16pt}}`
}

//GenIndex returns the default html file as a string.
func (d *DefaultBlog) GenIndex() string {
	return `<!DOCTYPE html>
	<html lang="en">
	
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta name="keywords" content="{{range .Keywords}} {{.}}, {{end}}">
		<link rel="stylesheet" href="index.css">
		<title>{{ .Title }}</title>
	</head>
	
	<body>
		{{range .Articles}}
		{{if .Public}}
		<article>
			<h1>{{.Title}}</h1>
			<h2>{{.Subtitle}}</h2>
			<hr>
			<div class="content">
				{{.HTML}}
			</div>
		</article>
		{{end}}
		{{end}}
	</body>
	
	</html>`
}


type Article struct {
	Public   bool     `json:"public"`
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle"`
	Tags     []string `json:"tags"`
	Body     string   `json:"body"`
	Slug     string   `json:"slug"`
	HTML 	template.HTML `json:"html"`
}
