package main

import(
	"github.com/hermannfass/textbase/textutil"
	"fmt"
)

func main() {
   r := textutil.Wrap("Das hübsche ß ist das schärfste" +
                      "S im Sprachraum. A B C D E F G H I", 8)
   fmt.Println(r)
}

