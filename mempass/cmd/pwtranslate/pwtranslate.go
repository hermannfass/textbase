package main

import(
	"github.com/hermannfass/textbase/mempass"
	"fmt"
	"strings"
	"os"
	"os/user"
	"path/filepath"
	"flag"
	"log"
)

// Translates password based on the keys given in secretsPath
// Command line argument: Password to translate to a code
// and/or optional flags.
func main() {
	// Default configuration path:
	user, uerr := user.Current()
	if uerr != nil {
		log.Fatal(uerr.Error())
	}
	defConfDir      := filepath.Join(user.HomeDir, "mempass")
	defSecretsPath  := filepath.Join(defConfDir, "mempass-secrets.txt")
	secretsPathFlag := flag.String("s", defSecretsPath,
	                   	`Path to the mempass "secrets" file`)
	inPathFlag      := flag.String("f", "",
	                   	"Path to input file to be coded")
	outFileFlag     := flag.String("o", "",
	                   	"Output file with encoded (paraphrased) passwords")
	flag.Usage = func() {
		fmt.Println(usageText())
		fmt.Println("OPTIONS / PARAMETERS")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Get secret descriptions for all characters:
	secretsPath := *secretsPathFlag
	fmt.Printf("Using configuration file %s\n", secretsPath)
	secrets := mempass.SecretsFromFile(secretsPath)

	// Prepare file handle for output:
	var outFH *os.File
	if *outFileFlag == "" {
		outFH = os.Stdout
	} else {
		var ferr error // Interface error
		outFH, ferr = os.Create(*outFileFlag)
		if ferr != nil {
			log.Fatal(ferr)
		}
		defer outFH.Close()
	}

	var result string
	if (*inPathFlag == "") {
		// No file specified => encode just one word
		var w string
		if (len(flag.Args()) == 0) {
			// Word not an argument => ask for it
			fmt.Print("Word to spell with secrets: ")
			_, rerr := fmt.Scanln(&w)
			if (rerr != nil) {
				log.Fatal(rerr)
			}				
		} else {
			w = flag.Args()[0]
		}
		fmt.Printf("The secrets spelling the word \"%s\" are:\n", w)
		codes := mempass.CodesForWord(w, secrets)
		result = strings.Join(codes, "\n") + "\n"
	} else {
		fmt.Printf("Translating file %s.\n", *inPathFlag)
		result = mempass.EncodeFile(*inPathFlag, secrets)
	}
	n, werr := outFH.WriteString(result)
	if werr != nil {
		fmt.Println("Error writing to file:", werr)
		fmt.Printf("Error type: %T\n", werr)
	} else {
		if (*inPathFlag != "") {
			fmt.Printf("%d bytes written to: %s\n", n, outFH.Name())
		}
	}
}

func usageText() string{
	return `
┌────────────────────────────────────────────────────┐
│                    pwtranslate                     │
│  plain text passwords that only you can interpret  │
└────────────────────────────────────────────────────┘
USAGE
  Translates unencrypted passwords to plain text
  that only the originator understands.

CONFIGURATION (Mempass Secrets File)
  The codes for describing each character of a password are defined
  in a configuration file, herein called the Mempass Secrets File.
  You can point on a specific one with thespecified with the -s option
  or left to the default.
  (Call mempass -h for more information on parameters.)
  Format:
  <Character>": "<Description>
  Each line shows a character followed by a colon (:), a space ( )
  and a secret that will tell only the originator which character
  this secret can get resolved to.
  Example:
    a: The 5th letter in the lyrics of Trixy's favourite song.
  Ok, make sure nobody knows who you call Trixy...
  Indeed, composing a Mempass Secrets File is not that easy. You
  need to be creative and find descriptions for all letters,
  digits, and special characters used in your passwords.
  Warnings regarding your Mempass Secrets File:
  - The description of a character should not tell what type of
    character it is. That means it should not start with "the
    number ..." or "special character ...".
  - Don't leave the unencrypted file on your computer.
    Once you have created the clear-text password file based on its
    descriptions you will not need it, thus: Encrypt it and put it
    some media that you keep at home.
`
}

