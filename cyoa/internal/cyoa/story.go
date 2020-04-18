package cyoa

import (
	"html/template"
	"log"
	"net/http"
)

//Story ...
type Story map[string]Chapter

//Chapter ...
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

//Option ...
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

//NewChapterHTTPHandler ...
func NewChapterHTTPHandler(story Story, defaultChapter string,
	tpl *template.Template, log *log.Logger) http.Handler {
	h := handler{story, defaultChapter, tpl, log}
	return h
}

type handler struct {
	story          Story
	defaultChapter string
	tpl            *template.Template
	log            *log.Logger
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path[1:]
	log.Printf("Requesting: %v\n", path)

	if path == "" {
		path = h.defaultChapter
	}

	h.tpl.Execute(w, h.story[path])
}
