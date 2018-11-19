package jwpurl

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/akutz/sortfold"
)

type JwpError struct {
	Status  string
	Message string
	Code    string
	Title   string
}
type JwpSuccess struct {
	Status string
	Tags   []Tags
}
type Tags struct {
	Playlist int64
	Name     string
	Videos   int64
}

func apiCall(u *url.URL) []string {
	var tags = []string{}
	response, err := http.Get(u.String())
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		var jwp JwpSuccess
		var jwpE JwpError
		if response.StatusCode != 200 {
			json.Unmarshal(contents, &jwpE)
			fmt.Println(jwpE)
		} else {
			json.Unmarshal(contents, &jwp)
			for _, tag := range jwp.Tags {
				tags = append(tags, tag.Name)
			}
			sortfold.Strings(tags)
		}
	}
	return tags
}

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}

//TagManager creates valid JWP management api url and makes the api call
//create and delete method take single string in a slice as parameter, update takes old name and new name
func TagManager(parameters []string, method string) []string {
	var tags = []string{}
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")
	apiURL := fmt.Sprintf("https://api.jwplatform.com/v1/accounts/tags/%v?", method)
	var URL *url.URL
	URL, err := url.Parse(apiURL)
	if err != nil {
		panic("boom")
	}

	randInt := rangeIn(10000000, 99999999)
	nonce := strconv.Itoa(randInt)
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	value := parameters
	paramas := url.Values{}
	paramas.Add("api_format", "json")
	paramas.Add("api_key", apiKey)
	paramas.Add("api_nonce", nonce)
	paramas.Add("api_timestamp", timestamp)
	signature := fmt.Sprintf("api_format=%s&api_key=%s&api_nonce=%s&api_timestamp=%s%s", "json", apiKey, nonce, timestamp, apiSecret)
	switch method {
	case "create":
		paramas.Add("name", value[0])
		signature = fmt.Sprintf("api_format=%s&api_key=%s&api_nonce=%s&api_timestamp=%s&name=%s%s", "json", apiKey, nonce, timestamp, url.PathEscape(value[0]), apiSecret)
	case "update":
		paramas.Add("name", value[0])
		paramas.Add("new_name", value[1])
		signature = fmt.Sprintf("api_format=%s&api_key=%s&api_nonce=%s&api_timestamp=%s&name=%s&new_name=%s%s", "json", apiKey, nonce, timestamp, url.PathEscape(value[0]), url.PathEscape(value[1]), apiSecret)

	case "delete":
		paramas.Add("name", value[0])
		signature = fmt.Sprintf("api_format=%s&api_key=%s&api_nonce=%s&api_timestamp=%s&name=%s%s", "json", apiKey, nonce, timestamp, url.PathEscape(value[0]), apiSecret)
	}
	URL.RawQuery = paramas.Encode()
	h := sha1.New()
	h.Write([]byte(signature))
	hash := h.Sum(nil)
	hash1 := fmt.Sprintf("%x", hash)
	paramas.Add("api_signature", hash1)
	URL.RawQuery = paramas.Encode()
	if method == "list" {
		tags = apiCall(URL)
		return tags
	}
	apiCall(URL)
	return nil
}
