package textutil

import(
	"unicode/utf8"
	"strings"
	"fmt"
)

func Wrap(s string, width int) string {
	lines := WrappedLines(s, width)
	var r string
	for _, l := range(lines) {
		r += strings.Join(l, " ") + "\n"
	}
	return r
}

func WrappedLines(s string, width int) [][]string {
	// To do: Consider to replace all \n and \t bytes by spaces.
	words := strings.Fields(s)
	var r [][]string  // Result
	var l []string    // Current line
	for i, w := range(words) {
		if complen(l, w) <= width {   // word still fits
			l = append(l, w)
		} else {                      // word is too long
			r = append(r, l)
			l = []string{w}
		}
		if i == len(words) - 1 {      // but the last word
			r = append(r, l)
		}
	}
	return r
}
			
func complen(words []string, newword string) int {
	return utf8.RuneCountInString(strings.Join(words, " ")) + 1 +
	       utf8.RuneCountInString(newword)
}

func main() {
	r := WrappedLines("Das hübsche ß ist das schärfste S im Sprachraum.", 9)
	fmt.Println(r)
	fmt.Println("----------- result ------------")
	for _, l := range(r) {
		fmt.Println(strings.Join(l, " "))
	}
}


/*

|-------------|
eins zwei drei
vier fünf sechs
sieben acht
neun zehn



*/

