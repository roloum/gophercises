package link

import (
	"reflect"
	"strings"
	"testing"
)

func TestOneLink(t *testing.T) {
	expected := []Link{{"http://google.com", "Google"}}
	s := "<div><p>xxx</p><p><a href=\"http://google.com\" rel=\"nofollow\" target=\"_blank\">Google</a></p></div>"
	r := strings.NewReader(s)
	links, _ := Parse(r)

	if !reflect.DeepEqual(links, expected) {
		t.Errorf("Expected: %v, received: %v", expected, links)
	}

}

func TestOneLinkExtraTags(t *testing.T) {
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
	links, _ := Parse(r)

	if !reflect.DeepEqual(links, expected) {
		t.Errorf("Expected: %v, received: %v", expected, links)
	}

}
