package server

import (
	"log"
	"net/http"

	"src/shell"
)

var appNavLinks = []shell.NavLink{
	{Location: "/", Text: "Home"},
	{Location: "/health", Text: "Health"},
	{Location: "/testPage", Text: "Test Page"},
}

type appSideBarData struct {
	Heading string
}

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	err := shell.Initialize(mux)
	if err != nil {
		panic("There was an error!")
	}

	// set up all routes
	mux.HandleFunc("/health", handleHealth)

	mux.HandleFunc("/testPage", shell.BuildShell(shell.Shell[testPageData, appSideBarData]{
		HeaderData: shell.HeaderData{
			Title:    "Test Page",
			PageName: "Test Page",
			Nav:      &appNavLinks,
		},
		SideBar: shell.NewStaticPagePart[appSideBarData]("AppSideBar", appSideBarData{Heading: "Test Side Bar"}),
		Body:    shell.NewPagePart[testPageData]("TestPageTemplate", getTestPageData),
	}))

	return mux
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	log.Printf("Health endpoint accessed")
	s := shell.Shell[string, string]{
		HeaderData: shell.HeaderData{
			Title:    "Health Check",
			PageName: "Health Check",
		},
	}

	shell.BuildShell(s)(w, r)
}

type testPageData struct {
	Number int64
}

func getTestPageData(*http.Request) (value testPageData) {
	value.Number = 42
	return
}
