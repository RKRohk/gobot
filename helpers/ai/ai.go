package ai

import (
	context "context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var cc *grpc.ClientConn
var err error

const (
	AI_LOCAL_URI = "localhost:50051"
	AI_PROD_URI  = "chatbotapi:50051"
)

var opts = grpc.WithInsecure()

var URI string

var prod = os.Getenv("PROD")

var client MessageServiceClient

func init() {
	log.Println(prod)
	if prod == "TRUE" {
		URI = AI_PROD_URI
	} else {
		URI = AI_LOCAL_URI
	}

	log.Println("URL selected is", URI)

	cc, err = grpc.Dial(URI, opts)

	if err != nil {
		log.Println("Error connecting to chatbotapi")
	} else {
		log.Println("Connected to chatbot api")
		client = NewMessageServiceClient(cc)

	}

}

func GetMessageResponse(message *tgbotapi.Message) (*MessageResponse, error) {
	ctx := context.TODO()

	messageConfig := &MessageRequest{Message: message.Text}
	return client.GetResponse(ctx, messageConfig)

}

func SendMessage(message string) (*emptypb.Empty, error) {
	ctx := context.TODO()

	messageConfig := &MessageRequest{Message: message}

	return client.Train(ctx, messageConfig)

}

func Close() {
	cc.Close()
}
