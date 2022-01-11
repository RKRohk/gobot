package search

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var cfg elasticsearch7.Config
var err error
var es7 *elasticsearch7.Client
var logger *log.Logger
var elasticSearchURL string = os.Getenv("ELASTICSEARCH_URL")
var bonsaiURL string = os.Getenv("BONSAI_URL")

//Document represents an elasticsearch document
type Document struct {
	FileID   string `json:"fileID,omitempty"`
	Data     string `json:"data,omitempty"`
	FileName string `json:"fileName"`
}

func init() {
	if len(elasticSearchURL) == 0 {
		elasticSearchURL = bonsaiURL
	}
	logger = log.New(os.Stdout, "helpers/search/ ", log.LstdFlags)

	cfg = elasticsearch7.Config{Addresses: []string{elasticSearchURL}}
	es7, err = elasticsearch7.NewClient(cfg)
	for err != nil {
		fmt.Printf("Error initializing elasticsearch client %v\n", err)
		fmt.Println("Trying in 5 seconds")
		time.Sleep(5 * time.Second)
		es7, err = elasticsearch7.NewClient(cfg)

	}

	_, err := es7.Ping()
	for err != nil {
		fmt.Printf("Error pinging elasticsearch %v\n", err)
		fmt.Println("Trying again in 5 seconds")
		time.Sleep(5 * time.Second)
		_, err = es7.Ping()
	}

	fmt.Println("Elasticsearch initialized")

}

//Index function takes a telagram document, converts it to base64 and sends it to elasticsearch to index
func Index(link string, repliedToDocument *tgbotapi.Document, hashTag string) {

	//Downloading the file
	res, err := http.Get(link)
	if err != nil {
		logger.Println("index.go: Error downloading file from url", err)
	} else {
		logger.Println(res.Body)
	}

	//Extracting the tag
	tag := strings.Replace(hashTag, "#", "", 1)
	fileName := repliedToDocument.FileName

	if true {
		file, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Println("index.go Unable to read downloaded file")
			return
		}

		//Encoding downloaded file to base64
		data := base64.StdEncoding.EncodeToString(file)
		elasticDoc := &Document{FileID: repliedToDocument.FileID, Data: data, FileName: fileName}

		body, err := json.Marshal(elasticDoc)

		index := esapi.IndexRequest{Index: tag, DocumentID: repliedToDocument.FileID, DocumentType: "_doc", Body: bytes.NewReader(body), Pipeline: "attachment"}

		ctx := context.Background()
		if res, err := index.Do(ctx, es7); err != nil {
			log.Println("index.go Error: ", err)
		} else {
			log.Println("index.go Created Index: ", res)
		}
	}

}

//Index function takes a telagram document, converts it to base64 and sends it to elasticsearch to index
func IndexBulk(link string, fileName string, fileID string, hashTag string) {

	//Downloading the file
	res, err := http.Get(link)
	if err != nil {
		logger.Println("index.go: Error downloading file from url", err)
	} else {
		logger.Println(res.Body)
	}

	//Extracting the tag
	tag := strings.Replace(hashTag, "#", "", 1)

	if true {
		file, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Println("index.go Unable to read downloaded file")
			return
		}

		//Encoding downloaded file to base64
		data := base64.StdEncoding.EncodeToString(file)
		elasticDoc := &Document{FileID: fileID, Data: data, FileName: fileName}

		body, err := json.Marshal(elasticDoc)

		index := esapi.IndexRequest{Index: tag, DocumentID: fileID, DocumentType: "_doc", Body: bytes.NewReader(body), Pipeline: "attachment"}

		ctx := context.Background()
		if res, err := index.Do(ctx, es7); err != nil {
			log.Println("index.go Error: ", err)
		} else {
			log.Println("index.go Created Index: ", res)
		}
	}

}
