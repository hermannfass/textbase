package pwsafetools

import(
	"fmt"
	"log"
	"os"
	"bufio"
	"strings"
	"errors"
)

type pwrec struct {
	title string
	user string
	password string
	url string
	email string
}

// Convert a pwSafe export file into a file ready to be used by mempass.
// Returns an error if not successful (i.e. error occurs) otherwise nil.
func ConvPwsafeExportFile(inpath string, outpath string) error {
	recs, rErr := ReadPwsafeExportFile(inpath)
	if rErr != nil {
		return rErr
	}
	var allStrs []string
   for _, r := range recs { // Construct string for each record
		// Add to the slice of strings
		rstr := r.title
		if r.url != "" {
			rstr += fmt.Sprintf(" (%s)", r.url)
		}
      if r.email != "" && r.email != r.user {
			rstr += fmt.Sprintf("\nEmail:    %s", r.email)
		}
		rstr += fmt.Sprintf("\nBenutzer: %s", r.user)
		rstr += fmt.Sprintf("\nPasswort: %s", r.password)
   	allStrs = append(allStrs, rstr) // Add record's string to slice allStrs
	}
   // Concatenate strings representing all records to a string
	s := strings.Join(allStrs, "\n\n") + "\n"
	err := os.WriteFile(outpath, []byte(s), 0600); 
	if err != nil {
		return err
	}
	fmt.Println("Input for Mempass file generation ready at", outpath)
	return nil
}


// Read a tab-delimited "pwSafe" export file and
// return a slice of pwrec elements.
func ReadPwsafeExportFile(path string) ([]pwrec, error) { 
	f, ferr := os.Open(path)
	if ferr != nil {
		return nil, ferr
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	// result := ""
	var result []pwrec
	for scanner.Scan() {
		l := scanner.Text()
		rec, err := convPwsafeExportLine(l)
		if err != nil {
			return nil, err
		} else {
			result = append(result, rec)
		}
	}
	return result, nil
}

// Convert a line from a "pwSafe" export file to a pwrec struct.
func convPwsafeExportLine(l string) (pwrec, error) {
	f := strings.Split(l, "\t") // All fields in that line
	if len(f) < 12 {
		log.Print("Not a complete record:\n", l)
		return pwrec{}, errors.New("Not a proper record:\n" + l)
	} else {
		title := f[0]
		user  := f[1]
		pass  := f[2]
		url   := f[3]
		mail  := f[10]
		if (strings.Contains(title, ".")) {
			title = strings.Split(title, ".")[1]
		}
		rec := pwrec{title, user, pass, url, mail}
		return rec, nil
	}
}

