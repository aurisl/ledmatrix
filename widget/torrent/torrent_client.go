package torrent

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"time"
)

type (
	UTorrent struct {
		username  string
		password  string
		token     string
		cookieJar http.CookieJar
		url       string
	}

	List struct {
		Build    int     `json:"build"`
		Torrents []Entry `json:"torrents"`
	}

	Entry struct {
		Hash          string
		Status        int
		Name          string
		Size          int
		Progress      float64
		Downloaded    int
		Uploaded      int
		Ratio         int
		UploadSpeed   int
		DownloadSpeed int
		Eta           int
		Label         string
		Remaining     uint64
	}
)

func (entry *Entry) UnmarshalJSON(data []byte) error {
	var value []interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		log.Printf("Error whilde decoding %value\n", err)
		return err
	}

	entry.Hash = value[0].(string)
	entry.Status = int(value[1].(float64))
	entry.Name, _ = value[2].(string)
	entry.Size, _ = value[3].(int)
	entry.Progress = value[4].(float64)
	entry.Downloaded, _ = value[5].(int)
	entry.Uploaded, _ = value[5].(int)
	entry.Ratio, _ = value[6].(int)
	entry.UploadSpeed, _ = value[7].(int)
	entry.DownloadSpeed, _ = value[8].(int)
	entry.Eta, _ = value[9].(int)
	entry.Eta, _ = value[10].(int)
	entry.Label, _ = value[11].(string)
	entry.Remaining = uint64(value[18].(float64))

	return nil
}

func NewUTorrentClient(url string, username string, password string) (*UTorrent, error) {

	UTorrent := &UTorrent{username: username, password: password, url: url}
	body := UTorrent.makeRequest("token.html", "GET")

	if len(body) == 0 {
		return nil, errors.New("empty body response")
	}

	UTorrent.token = extractAuthToken(body)

	return UTorrent, nil
}

func (UTorrent *UTorrent) getList() (*List, error) {

	body := UTorrent.makeRequest("", "GET")

	if len(body) == 0 {
		return nil, errors.New("empty body response")
	}

	torrentList := &List{}
	err := json.Unmarshal(body, torrentList)

	if err != nil {
		log.Println(err.Error())
	}

	return torrentList, nil

}
func extractAuthToken(body []byte) string {

	regularExpression, _ := regexp.Compile(`<div id=(?:\'|")token(?:\'|")[^>]+>(.*)</div>`)

	matches := regularExpression.FindAllSubmatch(body, -1)
	token := string(matches[0][1])

	return token
}

func (UTorrent *UTorrent) makeRequest(path string, method string) []byte {

	req, _ := http.NewRequest(method, UTorrent.url+path, nil)

	req.URL.User = url.UserPassword(UTorrent.username, UTorrent.password)

	q := req.URL.Query()
	q.Add("list", "1")
	q.Add("token", UTorrent.token)

	req.URL.RawQuery = q.Encode()

	if UTorrent.cookieJar == nil {
		cookie, _ := cookiejar.New(nil)
		UTorrent.cookieJar = cookie
	}

	client := &http.Client{
		Jar:     UTorrent.cookieJar,
		Timeout: time.Second * 3,
	}

	response, err := client.Get(req.URL.String())

	if err != nil {
		log.Println(err.Error())
		return []byte{}
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Println(err.Error())
		return []byte{}
	}

	err = response.Body.Close()
	if err != nil {
		log.Printf("An err occurred while closing response body '%s'", err.Error())
		return nil
	}

	return body
}
