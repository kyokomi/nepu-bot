package twitter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/shogo82148/go-shuffle"
)

var postImages map[string]bool

type TwitterSearchResponse struct {
	Statuses []struct {
		Entities struct {
			Media []struct {
				MediaURL      string `json:"media_url"`
				MediaURLHTTPs string `json:"media_url_https"`
			} `json:"media"`
		} `json:"entities"`
	} `json:"statuses"`
	ExpireTime time.Time
}

func (r TwitterSearchResponse) GetMediaURLs() []string {
	urls := map[string]bool{}
	for _, s := range r.Statuses {
		for _, m := range s.Entities.Media {
			if postImages[m.MediaURL] {
				continue
			}
			if urls[m.MediaURL] {
				continue
			}
			urls[m.MediaURL] = true
		}
	}

	results := make([]string, len(urls))
	i := 0
	for url, _ := range urls {
		results[i] = url
		i++
	}
	return results
}

func (r TwitterSearchResponse) Expire() bool {
	return r.ExpireTime.Before(time.Now())
}

func newAccessToken(key, secret string) (string, error) {
	if key == "" {
		key = os.Getenv("TWITTER_API_KEY")
	}
	if secret == "" {
		secret = os.Getenv("TWITTER_API_SECRET")
	}
	token := base64.StdEncoding.EncodeToString([]byte(url.QueryEscape(key) + ":" + url.QueryEscape(secret)))

	client := &http.Client{}
	data := url.Values{"grant_type": {"client_credentials"}}

	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Authorization", "Basic "+token)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var f interface{}
	json.Unmarshal(body, &f)
	m := f.(map[string]interface{})

	return m["access_token"].(string), nil
}

const (
	searchAPIURL = "https://api.twitter.com/1.1/search/tweets.json"
)

var responseCache TwitterSearchResponse

func searchImages(accessToken, keyword string, count int) ([]string, error) {
	urls := responseCache.GetMediaURLs()
	if len(urls) == 0 || responseCache.Expire() {
		apiURL := fmt.Sprintf("%s?q=%s&count=%d&include_entities=true",
			searchAPIURL,
			url.QueryEscape("filter:images "+keyword),
			100,
		)
		searchReq, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			return nil, err
		}
		searchReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		searchReq.Header.Set("Authorization", "Bearer "+accessToken)

		client := &http.Client{}
		searchRes, err := client.Do(searchReq)
		if err != nil {
			return nil, err
		}
		defer searchRes.Body.Close()

		searchBody, err := ioutil.ReadAll(searchRes.Body)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(searchBody, &responseCache); err != nil {
			return nil, err
		}

		urls = responseCache.GetMediaURLs()
		if len(urls) == 0 {
			urls, _ = searchImages(accessToken, "ネプテューヌ", count)
		}
	}

	shuffle.Strings(urls)

	if postImages == nil {
		postImages = make(map[string]bool, 0)
	}
	if len(urls) < count {
		count = len(urls)
	}
	for _, imageURL := range urls[:count] {
		postImages[imageURL] = true
	}

	return urls[:count], nil
}
