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
)

func AccessToken() error {
	key := os.Getenv("TWITTER_API_KEY")
	secret := os.Getenv("TWITTER_API_SECRET")
	token := base64.StdEncoding.EncodeToString([]byte(url.QueryEscape(key) + ":" + url.QueryEscape(secret)))

	client := &http.Client{}
	data := url.Values{"grant_type": {"client_credentials"}}

	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Authorization", "Basic "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var f interface{}
	json.Unmarshal(body, &f)
	m := f.(map[string]interface{})

	fmt.Println(m["access_token"].(string))

	searchReq, err := http.NewRequest("GET", "https://api.twitter.com/1.1/search/tweets.json?q="+url.QueryEscape("イストワール"), nil)
	if err != nil {
		return err
	}
	searchReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	searchReq.Header.Set("Authorization", "Bearer "+m["access_token"].(string))

	searchRes, err := client.Do(searchReq)
	if err != nil {
		return err
	}
	defer searchRes.Body.Close()

	fmt.Println(searchRes.Status)

	searchBody, err := ioutil.ReadAll(searchRes.Body)
	if err != nil {
		return err
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(searchBody, &jsonMap); err != nil {
		return err
	}
	fmt.Println(jsonMap)

	return nil
}
