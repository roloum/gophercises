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
