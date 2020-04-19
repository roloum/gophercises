package cyoa

import (
	"html/template"
	"log"
	"net/http"
)

//Stories ...
type Stories map[string]Story

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

type handler struct {
	story          Story
	defaultChapter string
	log            *log.Logger
	tpl            *template.Template
}

//Default HTML template
var defaultTemplate = `Hello!`

//WithNewTemplate option for setting up new template
func WithNewTemplate(tpl *template.Template) func(h *handler) {
	return func(h *handler) {
		h.tpl = tpl
	}
}

//NewChapterHTTPHandler ...
func NewChapterHTTPHandler(story Story, defaultChapter string,
	log *log.Logger, options ...func(h *handler)) http.Handler {

	//Create template for default HTML
	tpl := template.Must(template.New("").Parse(defaultTemplate))
	//Setting up handler with default options
	h := handler{story, defaultChapter, log, tpl}

	//Applying Options
	for _, opt := range options {
		opt(&h)
	}

	return h
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path[1:]
	log.Printf("Requesting: %v\n", path)

	if path == "" {
		path = h.defaultChapter
	} else if _, ok := h.story[path]; !ok {
		log.Printf("Story not found: %v\n", path)
		http.Error(w, "Story not found", http.StatusNotFound)
		return
	}

	if err := h.tpl.Execute(w, h.story[path]); err != nil {
		h.log.Printf("Error executing template: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
