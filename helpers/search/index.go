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

type Document struct {
	FileID   string `json:"fileID,omitempty"`
	Data     string `json:"data,omitempty"`
	FileName string `json:"fileName"`
}

func init() {
	cfg = elasticsearch7.Config{Addresses: []string{"http://52.172.252.187:9200"}}
	es7, err = elasticsearch7.NewClient(cfg)
	logger = log.New(os.Stdout, "helpers/search/index.go: ", log.LstdFlags)
	if err != nil {
		log.Panic(err)
	} else {
		log.Println(es7)
	}
}

func Index(link string, repliedToDocument *tgbotapi.Document, hashTag string) {
	res, err := http.Get(link)
	if err != nil {
		logger.Println("index.go: Error downloading file from url", err)
	} else {
		logger.Println(res.Body)
	}
	tag := strings.Replace(hashTag, "#", "", 1)
	fileName := repliedToDocument.FileName
	if true {
		file, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Println("Unable to read downloaded file")
		}
		data := base64.StdEncoding.EncodeToString(file)
		elasticDoc := &Document{FileID: repliedToDocument.FileID, Data: data, FileName: fileName}

		body, err := json.Marshal(elasticDoc)

		index := esapi.IndexRequest{Index: tag, DocumentType: "_doc", Body: bytes.NewReader(body), Pipeline: "attachment"}

		ctx := context.Background()
		if res, err := index.Do(ctx, es7); err != nil {
			log.Println("Error: ", err)
		} else {
			log.Println("Created Index: ", res)
		}
	}

}
