package main

import "regexp"

func main() {

}

func normalize(phone string) string {
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(phone, "")
}
