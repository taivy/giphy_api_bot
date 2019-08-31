package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func queryGif(query string, resLimit int) []string {
	ApiKey := "IHuGuz4N236oJuDLKDoYpG4ckGqz35D2"
	URL := "http://api.giphy.com/v1/gifs/search"

	req, _ := http.NewRequest("GET", URL, nil)
	q := req.URL.Query()
	q.Add("api_key", ApiKey)
	q.Add("q", query)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	map_ := map[string]interface{}{}
	json.Unmarshal(body, &map_)
	data := map_["data"]
	gifs := data.([]interface{})
	i := 0
	results := make([]string, resLimit)
	for i < resLimit {
		entry := gifs[i]
		gifUrl := entry.(map[string]interface{})["url"].(string)
		i += 1
		results = append(results, gifUrl)
	}
	fmt.Println(results)
	return results
}

func GetGif(input string) (string, []string) {
	split := strings.Split(input, " ")
	var err_reply string
	resLimit := 5
	if len(split) < 2 {
		err_reply =
			`Incorrect command format`
	} else {
		var err error
		resLimit, err = strconv.Atoi(split[len(split)-1])
		if err != nil {
			err_reply =
				`Incorrect format. Third argument is not integer. Command must have format:
			/get_gif query [result number limit]
			result number limit is optional and defaults to 5. You can set it to value from 0 to 25
			`
		}
		if resLimit > 25 {
			err_reply =
				"Too many results requests. Limit is 25"
		}
	}
	var res []string
	if err_reply == "" {
		query := strings.Join(split[2:], " ")
		res = queryGif(query, resLimit)
	}
	return err_reply, res

}
