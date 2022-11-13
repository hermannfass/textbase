package mempass;

import(
	"fmt"
	"strings"
	"log"
	"regexp"
	"os"
	"bufio"
	"io/ioutil"
)

// Reads the users configuration file with one character and its
// "secret code" (i.e. paraphrase) per line, separated by ": ".
// The return value is a map with the individual characters as
// keys and their respective code as values.
func SecretsFromFile(path string) map[string]string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("The file", path, "does not exist.")
		log.Fatal(err)
	}
	defer f.Close()
	fs := bufio.NewScanner(f)
	// The following should be redundant as bufio.ScanLines is the
	// pre-defined implementation of type SplitFunc. 
	// fs.Split(bufio.ScanLines) // set split function for scanner.
	secrets := map[string]string{}
	for fs.Scan() {
		// To do: Separator in next line should perhaps be configurable?
		pair := strings.SplitN(fs.Text(), ": ", 2) 
		if len(pair) == 2 {
			secrets[pair[0]] = pair[1]
		} else {
			secrets[pair[0]] = pair[0]
		}
	}
	return secrets
}

// Read a file, replace the password xyz in "Password: xyz" with an
// indented set of lines paraphrasing the password.
func EncodeFile(path string, secrets map[string]string) string {
	b, err := ioutil.ReadFile(path)  // reads a Slice of bytes ([]byte)
	result := ""
	if err != nil {
		fmt.Println("The file", path, "does not exist.")
		log.Fatal(err)
	}
	s := string(b)  // b is a []byte
	pwLabelPat := "Passwor(?:d|t):" // To do: Make configurable?
	// Pattern to detect a password (non-spaces after PW label)
	rePw := regexp.MustCompile(`(?msi)` + pwLabelPat + `\s*(\S+)`)
	// Function to be used in ReplaceAllStringFunc for replacing
	// found passwords by the return value of this function (i.e.
	// the paraphrased password:
	rf := func(rawStr string) string {
		// Search for submatches (get label and password)
		pwm := rePw.FindStringSubmatch(rawStr)
		// Get the password as the 1st submatch (0 is whole match)
		pw := pwm[1]
		// Get a slice of strings describing the password chars
		codes := CodesForWord(pw, secrets)
		// Put the slices together as a new multi-line string:
		var codeList string
		for _, v := range codes {
			codeList += "  " + v + "\n" 
		}
		return codeList
	}
	result = rePw.ReplaceAllStringFunc(s, rf)
	return result
}

func EncodeWord(clear string, secrets map[string]string) string {
	codes := CodesForWord(clear, secrets)
	return strings.Join(codes, "\n")
}

func CodesForWord(clear string, secrets map[string]string) []string {
	codes := make([]string, 0) 
	for _, r := range clear {   // all runes; ignoring index
		if sec, ok := secrets[string(r)]; ok {
			codes = append(codes, sec)
		} else {
			codes = append(codes, string(r))
		}
	}
	return codes
}





