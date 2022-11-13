package main;

import(
	"github.com/hermannfass/textbase/pwsafetools"
	"os/user"
	"path/filepath"
	"flag"
	"fmt"
	"log"
)

// Prepares a file for being used by mempass:
// Reads the pwSafe export file: ~/mempass/pwsafe-export.txt
// Extracts relevant information into: ~/mempass/mempass-pwlist.txt
// To do:
// Paths are hardcoded but should be relative to $HOME (cross platform)
// and ideallly also optional command line argument.
func main() {
	// Default file path:
	user, uerr := user.Current()
	if uerr != nil {
		log.Fatal(uerr.Error())
	}
	defPwsafePath := filepath.Join(user.HomeDir, "mempass", "pwsafe.txt")
	defResultPath := filepath.Join(user.HomeDir, "mempass", "pwlist.txt")
	pwsafePathFlag := flag.String("f", defPwsafePath, "Path of pwSafe export file")
	resultPathFlag := flag.String("o", defResultPath, "Output path")
	flag.Parse()
	fmt.Println("Converting pwSafe export file at", *pwsafePathFlag)
	fmt.Println("to an encoded password list at", *resultPathFlag)
	err := pwsafetools.ConvPwsafeExportFile(*pwsafePathFlag, *resultPathFlag)
	if err != nil {
		fmt.Printf("An error occured:\n%s\n", err.Error())
	}
}

