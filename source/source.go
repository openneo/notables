package source

import (
	"io/ioutil"
	"net/http"
	"regexp"
	. "github.com/openneo/neopets-notables-go/notables"
)

var (
	petNamePattern = regexp.MustCompile(
		`/petlookup\.phtml\?pet=([a-zA-Z0-9_]+)`)
	imageHashPattern = regexp.MustCompile(
		`http://pets\.neopets\.com/cp/([a-z0-9]+)/1/2\.png`)
)

func GetNotable(maxTries uint64) (Notable, bool) {
	for i := uint64(0); i < maxTries; i++ {
		notable := requestNotable()
		if notable.PetName != "" {
			return notable, true
		}
	}
	return Notable{}, false
}

func requestNotable() Notable {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://www.neopets.com", nil)
	if err != nil {
		return Notable{}
	}
	req.Header.Add("User-Agent", "OpenNeo Notables Tracker")

	resp, err := client.Do(req)
	if err != nil {
		return Notable{}
	}

	// CONSIDER: match on a buffer instead of loading in the whole string?
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Notable{}
	}
	resp.Body.Close()
	body := string(bodyBytes)

	petNameMatch := petNamePattern.FindStringSubmatch(body)
	if len(petNameMatch) < 2 {
		return Notable{}
	}

	imageHashMatch := imageHashPattern.FindStringSubmatch(body)
	if len(imageHashMatch) < 2 {
		return Notable{}
	}

	return Notable{petNameMatch[1], imageHashMatch[1]}
}
