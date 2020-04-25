package link

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestParseOneLink(t *testing.T) {
	expected := []Link{{"http://google.com", "Google"}}
	s := "<div><p>xxx</p><p><a href=\"http://google.com\" rel=\"nofollow\" target=\"_blank\">Google</a></p></div>"
	r := strings.NewReader(s)
	links, err := Parse(r)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(links, expected) {
		t.Errorf("Expected: %v, received: %v", expected, links)
	}

}

func TestParseOneLinkExtraTags(t *testing.T) {
	expected := []Link{
		{"http://google.com", "Output Google Maps"},
		{"http://yahoo.com", "Yahoo!"},
	}
	s := `
		<div>
			<p>
				xxx
			</p>
			<p>
				<a href="http://google.com">Output <b>Google</b> <i>Maps</i></a>
				some more text
				<a href="http://yahoo.com">Yahoo!</a>
			</p>
		</div>`
	r := strings.NewReader(s)
	links, err := Parse(r)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(links, expected) {
		t.Errorf("Expected: %v, received: %v", expected, links)
	}

}

func TestParseEx1(t *testing.T) {
	expected := []Link{
		{"/other-page", "A link to another page"},
	}
	err := testParseFile("tests/ex1.html", expected)
	if err != nil {
		t.Error(err)
	}
}

func TestParseEx2(t *testing.T) {
	expected := []Link{
		{"https://www.twitter.com/joncalhoun", "Check me out on twitter"},
		{"https://github.com/gophercises", "Gophercises is on Github!"},
	}
	err := testParseFile("tests/ex2.html", expected)
	if err != nil {
		t.Error(err)
	}
}

func TestParseEx3(t *testing.T) {
	expected := []Link{
		{"#", "Login"},
		{"/lost", "Lost? Need help?"},
		{"https://twitter.com/marcusolsson", "@marcusolsson"},
	}
	err := testParseFile("tests/ex3.html", expected)
	if err != nil {
		t.Error(err)
	}
}

func TestParseEx4(t *testing.T) {
	expected := []Link{
		{"/dog-cat", "dog cat"},
	}
	err := testParseFile("tests/ex4.html", expected)
	if err != nil {
		t.Error(err)
	}
}

func testParseFile(fileName string, expected []Link) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	links, err := Parse(file)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(links, expected) {
		return fmt.Errorf("Expected: %v, received: %v", expected, links)
	}
	return nil
}
