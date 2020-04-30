package sitemap

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

var serverName = "http://localhost:8001"

func TestGetPages(t *testing.T) {

	serverChan := make(chan int)
	go func() {
		startUpServer(serverChan, t)
	}()

	//t.Log("waiting for server to start...")
	<-serverChan
	//t.Log("testing... wait another sec...")
	//time.Sleep(time.Second)
	//t.Log("done testing...")

	l := log.New(os.Stdout, "sitemap.test", log.Ldate)
	pages, err := GetPages(serverName, 0, l)
	if err != nil {
		t.Error(err)
	}

	expected := []string{
		fmt.Sprintf("%s/", serverName),
		fmt.Sprintf("%s/page1", serverName),
		fmt.Sprintf("%s/page2", serverName),
		fmt.Sprintf("%s/page3", serverName),
		fmt.Sprintf("%s/page4", serverName),
	}

	if !reflect.DeepEqual(expected, pages) {
		t.Errorf("Expected: %v, received: %v", expected, pages)
	}
}

func TestGetPagesDepth(t *testing.T) {

	l := log.New(os.Stdout, "sitemap.test", log.Ldate)
	pages, err := GetPages(serverName, 2, l)
	if err != nil {
		t.Error(err)
	}

	expected := []string{
		fmt.Sprintf("%s/", serverName),
		fmt.Sprintf("%s/page1", serverName),
		fmt.Sprintf("%s/page2", serverName),
	}

	if !reflect.DeepEqual(expected, pages) {
		t.Errorf("Expected: %v, received: %v", expected, pages)
	}
}

func startUpServer(serverChan chan int, t *testing.T) {

	sm := http.NewServeMux()
	sm.HandleFunc("/", home)
	sm.HandleFunc("/page1", page1)
	sm.HandleFunc("/page2", page2)
	sm.HandleFunc("/page3", page3)
	sm.HandleFunc("/page4", page4)

	server := &http.Server{
		Addr:    ":8001",
		Handler: sm,
	}

	//t.Log("start up server...")

	go func() {
		if err := server.ListenAndServe(); err != nil {
			t.Fatal(err)
		}
	}()

	//t.Log("waiting...")
	time.Sleep(time.Second)
	//t.Log("done... send server signal...")
	serverChan <- 1
	//t.Log("sent server signal...")

}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w,
		`<a href="/page1">Page 1</a><br>
<a href="https://github.com">Github</a>`)
}

func page1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<a href="%s/page2">Page 2</a>`, serverName)
}

func page2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<a href="/page3">Page 3</a>`)
}

func page3(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<a href="/page4">Page 4</a>`)
}

func page4(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<a href="/">Home</a>`)
}
