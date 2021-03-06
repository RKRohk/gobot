package search

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

//Search function takes a tag and a query string and returns an array of fileIDs which
//Contain that particular query
func Search(hashTag string, query string) []string {

	searchBody := `{
		"query": {
        	"match": {
            "attachment.content": "%s"
       		}
 		},
    	"fields": [
        	"fileID",
        	"fileName"
    	],
    	"_source": "false"
	}`

	searchBody = fmt.Sprintf(searchBody, query)

	sr := esapi.SearchRequest{Index: []string{hashTag}, Body: strings.NewReader(searchBody), Pretty: true, Timeout: 100}
	ctx := context.Background()
	response, err := sr.Do(ctx, es7)
	if err != nil {
		logger.Println("Error occured while searching", err)
	} else {
		var res Result
		err := json.NewDecoder(response.Body).Decode(&res)
		if err != nil {
			logger.Println("Error", err)
		} else {
			logger.Println("Search response for ", hashTag, "\n", response)
			fileIDs := make([]string, 0)
			for _, hit := range res.Hits.Hits {
				for _, fileID := range hit.Fields.FileID {
					fileIDs = append(fileIDs, fileID)
				}
			}
			return fileIDs
		}
	}
	return []string{}

}
