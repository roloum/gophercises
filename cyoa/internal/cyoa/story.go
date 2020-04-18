package cyoa

import (
	"html/template"
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

//NewStoryHTTPHandler ...
func NewStoryHTTPHandler(story Story, tpl *template.Template) http.Handler {
	h := handler{story, tpl}
	return h

}

type handler struct {
	story Story
	tpl   *template.Template
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hello := "Hello world"
	w.Write([]byte(hello))
}
