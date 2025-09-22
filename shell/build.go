package shell

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"

	"src/shell/assets"
)

var shellTemplateCache *template.Template
var appTemplateCache *template.Template

type DataFetcher[T any] func(*http.Request) T

type PagePart[T any] struct {
	TemplateName string
	ProvideData  DataFetcher[T]
}

func NewPagePart[T any](templateName string, f DataFetcher[T]) *PagePart[T] {
	return &PagePart[T]{
		TemplateName: templateName,
		ProvideData:  f,
	}
}

func NewStaticPagePart[T any](templateName string, data T) *PagePart[T] {
	return &PagePart[T]{
		TemplateName: templateName,
		ProvideData: func(r *http.Request) T {
			log.Printf("Providing static data: %o", data)
			return data
		},
	}
}

func (part *PagePart[T]) Render(w http.ResponseWriter, r *http.Request) {
	log.Println("In Render function")
	log.Print("%o", part)
	var data T
	if part.ProvideData != nil {
		data = part.ProvideData(r)
	}
	appTemplateCache.ExecuteTemplate(w, part.TemplateName, data)
}

func (part *PagePart[T]) renderShellPart(w http.ResponseWriter, r *http.Request) {
	log.Println("In renderShellPart function")
	log.Print("%o", part)
	var data T
	if part.ProvideData != nil {
		data = part.ProvideData(r)
	}
	shellTemplateCache.ExecuteTemplate(w, part.TemplateName, data)
}

type NavLink struct {
	Location string
	Text     string
}

type Shell[T any, S any] struct {
	HeaderData
	SideBar *PagePart[S]
	Body    *PagePart[T]
}

type HeaderData struct {
	Title    string
	PageName string
	UserInfo *struct {
		name string
		id   string
	}
	Nav *[]NavLink
}

func BuildShell[T any, S any](shell Shell[T, S]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := NewStaticPagePart[HeaderData]("HEADER", HeaderData{
			shell.Title,
			shell.PageName,
			shell.UserInfo,
			shell.Nav,
		})
		header.renderShellPart(w, r)
		if shell.SideBar != nil {
			shell.SideBar.Render(w, r)
		}
		if shell.Body != nil {
			shell.Body.Render(w, r)
		}
		footer := NewStaticPagePart[any]("FOOTER", nil)
		footer.renderShellPart(w, r)
	}
}

func Initialize(mux *http.ServeMux) error {
	// List all files inside the embedded filesystem
	err := fs.WalkDir(assets.Static, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			log.Println("Directory:", path)
		} else {
			log.Println("File:", path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	shellTemplateCache = template.Must(template.ParseFS(templates, "templates/*.tmpl"))
	// shellTemplateCache = template.Must(template.ParseGlob("shell/templates/*.tmpl"))
	appTemplateCache = template.Must(template.ParseGlob("internal/templates/*.tmpl"))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServerFS(assets.Static)))

	return nil
}

func Test() {
	t, _ := template.New("t").Parse(`
	<html>
		<head>
			<title>Rendered Template</title>
		</head>
		<body>
			<h1>Rendered Template</h1>
		</body>
	</html>
	`)
	data := make(map[string]string)
	t.Execute(os.Stdout, data)
}
