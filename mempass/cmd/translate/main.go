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
	                               "Output file with encoded passwords")
	flag.Usage = func() {
		fmt.Println("Usage:\nTranslates unencrypted passwords by",
		            "paraphrases (codes) for each character.\n\n",
		            "The codes are defined in a configuration file\n",
		            "(mempass secrets file, option -s) in which\n",
		            "each line shows a character followed by \": \" and\n",
		            "then the description that will tell the originator\n",
		            "of these secrets what character it represents.\n\n",
		            "The word or text to be translated is specified as\n",
		            "command line argument or in an input file specified\n",
		            "in option -f. In that case each word prefixed by\n",
		            "\"Password:\" gets replaced by the respective codes.")
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
	// outFH.Sync()

}

