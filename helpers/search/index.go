package search

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var cfg elasticsearch7.Config
var err error
var es7 *elasticsearch7.Client
var logger *log.Logger
var elasticSearchURL string = os.Getenv("ELASTICSEARCH_URL")

//Document represents an elasticsearch document
type Document struct {
	FileID   string `json:"fileID,omitempty"`
	Data     string `json:"data,omitempty"`
	FileName string `json:"fileName"`
}

func init() {
	if len(elasticSearchURL) == 0 {
		elasticSearchURL = "http://search:9200"
	}
	cfg = elasticsearch7.Config{Addresses: []string{elasticSearchURL}}
	es7, err = elasticsearch7.NewClient(cfg)
	logger = log.New(os.Stdout, "helpers/search/ ", log.LstdFlags)
	if err != nil {
		log.Panic(err)
	} else {
		_, err := es7.Ping()
		if err != nil {
			log.Panicln("Error connecting to elasticsearch", err)
		} else {
			log.Println("ElasticSearch Initialized!")
		}

	}
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

		index := esapi.IndexRequest{Index: tag, DocumentType: "_doc", Body: bytes.NewReader(body), Pipeline: "attachment"}

		ctx := context.Background()
		if res, err := index.Do(ctx, es7); err != nil {
			log.Println("index.go Error: ", err)
		} else {
			log.Println("index.go Created Index: ", res)
		}
	}

}
