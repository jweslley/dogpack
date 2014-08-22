package main

import (
	"io"
	"log"
	"text/template"
)

type tmplData map[string]interface{}

var templates = make(map[string]*template.Template)

func init() {
	for name, html := range templateHTML {
		t := template.New(name) //.Funcs(templateFuncs)
		template.Must(t.Parse(html))
		template.Must(t.Parse(baseHTML))
		templates[name] = t
	}
}

func render(w io.Writer, name string, data tmplData) {
	err := templates[name].ExecuteTemplate(w, "root", data)
	if err != nil {
		log.Println(err)
	}
}

const baseHTML = `
{{ define "root" }}
<html>
  <head>
    <title>{{.Title}}</title>
    <style type="text/css">
      body {
        background-color: #fff;
        font-family: "Helvetica Neue",Helvetica,Arial,sans-serif;
        font-size: 14px;
      }
      #container {
        width: 90%;
        max-width: 750px;
        padding-right: 15px;
        padding-left: 15px;
        margin-right: auto;
        margin-left: auto;
      }
      ul {
        list-style: none;
        padding: 0;
        margin: 5px 0;
      }
      li { color: #333; }
      ul.servers > li {
        padding: 10px;
        margin: 5px;
        border: 1px solid #ccc;
      }
      ul.servers > li.green { border-left: 5px solid #1abc9c; }
      ul.servers > li.yellow { border-left: 5px solid #f1c40f; }
      ul.servers > li.red { border-left: 5px solid #e74c3c; }
      ul.servers > li.gray { border-left: 5px solid #bdc3c7; }
			a, a:visited {
				color: #0074D9;
        text-decoration: none;
        background: none repeat scroll 0 0 transparent;
      }
      li.title {
        font-size: 1.5em;
        font-weight: 500;
        line-height: 1.1;
        padding-bottom: 5px;
      }
      li.green li.title { color: #1abc9c; }
      li.yellow li.title { color: #f1c40f; }
      li.red li.title { color: #e74c3c; }
      li.gray li.title { color: #bdc3c7; }
    </style>
  </head>
  <body>
    <div id="container">
      <h1>{{.Title}}</h1>
			{{ template "body" . }}
    </div>
  </body>
</html>
{{ end }}
`

var templateHTML = map[string]string{
	"index": `
{{ define "body" }}
	<ul class="servers">
	{{range .Nodes}}
		<li class="red">
			<a href="/servers/{{.Server.Localhostname}}">
				<ul class="info">
					<li class="title">{{.Server.Localhostname}}</li>
					<li>{{.Platform.Memory}} - {{.IpAddress}}</li>
					<li>{{.Platform.Name}} {{.Platform.Release}} {{.Platform.Machine}}</li>
				</ul>
			</a>
		</li>
	{{end}}
	</ul>
{{ end }}
`,
	"server": `
{{ define "body" }}
	<ul class="servers">
		<li class="red">
			<ul class="info">
				<li class="title">{{.Monit.Server.Localhostname}}</li>
				<li>{{.Monit.Platform.Memory}} - {{.Monit.IpAddress}}</li>
				<li>{{.Monit.Platform.Name}} {{.Monit.Platform.Release}} {{.Monit.Platform.Machine}}</li>
			</ul>
		</li>
	</ul>
	<ul class="services">
	{{range .Monit.Service}}
		<li class="red">
			<ul class="info">
				<li class="title">{{.Name}}</li>
			</ul>
		</li>
	{{ end }}
	</ul>
{{ end }}
`,
}
