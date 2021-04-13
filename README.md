# Gobot

A modular telegram bot written using GoLang. Utilizing microservices architecture


## Technologies Used
 - GoLang
 - Docker
 - Python
 - MongoDB
 - gRPC
 - AWS (Deployed)


## About

This bot started as a small project to store class notes and keep assignment reminders on telegram. 
Slowly and gradually, a lot of features have been added to it


## Structure

```
                --- Bot (written on GoLang)
                |
                |
                --- AI gRPC endpoint (written on Python)
                |
                |
                --- MongoDB instance
                |
                |
 Docker Daemon  --- Mongo Express (to view the database)
                |
                |
                --- ElasticSearch1 (load balanced)
                |
                |
                --- ElasticSearch2 (load balanced)
  
```


## Features

### AI

The bot learns from the speaking habits of it's users and learns to talk back. It is surprisingly good, almost seems like a real person is talking back.
The AI feature is handled by a python backend. The bot communicates with it using gRPC. 

To use this feature, just start talking to the bot


### DadJoke

`/dadjoke` - Sends a random Dad joke to the chat

### Echo

`/echo Hello World`
replies:
`Hello World`

Echos back the message sent

### Reminder
`/remind 25/02/2000 12:00 IST Rohan's Birthday`
replies:
```Okay, I will remind you about
 Rohan's Birthday
 
 on 25/02/2000 12:00 AM
```

Keeps reminders

### Save
`/save #<tag>`
Saves a note with a tag. Also if a file is given, it is passed through a pipeline where OCR is run on it, 
it is indexed and ready to be searched, like Google! but for your files.

### Get
`/get #<tag>`
Gets all notes with the given tag


### Search
`/search #<tag> <query>` - returns the file which contains the notes on query. 
The query need not be exact, as the search engine uses fuzzy algorithms and replaces words with synonyms to find the perfect match. 

### Slap
`/slap Rohan`
replies:
`Rohan hit Rohan around a bit with a trout`

Slaps the person the command is used against

## StickerSlap
`/stickerslap`

Reply to a user with this command and the bot will slap them with a sticker


## Deployment

This bot can be deployed on any linux server with Docker installed. 
Using docker compose, it can also be deployed on services like AWS Elastic Beanstalk, DigitalOcean Application Service, etc

## Docker Compose Structure
```
version: "3"
services:
  bot:
    image: rkrohk/gobot:latest
    env_file: "./.env.gobot"
    restart: always

  mongo:
    image: mongo
    container_name: mongo
    ports:
      - 27017:27017
    expose:
      - 27017
    volumes:
      - mongodata:/data/db

  chatbotapi:
    image: rkrohk/chatbotapi
    container_name: chatbotapi
    hostname: chatbotapi
    expose:
      - 50051

    environment:
      - MONGO_URI=mongo

    volumes:
      - chatbotapi_data:/src/app/data
      - nltk_data:/root/nltk_data

  mongo-express:
    image: mongo-express
    environment:
      - ME_CONFIG_MONGODB_SERVER=mongo
      - ME_CONFIG_MONGODB_PORT=27017
    ports:
      - "8080:8081"


  es01:
    image: rkrohk/elasticsearch
    container_name: es01
    environment:
      - node.name=es01
      - cluster.name=es-docker-cluster
      - discovery.seed_hosts=es02
      - cluster.initial_master_nodes=es01,es02  
      - bootstrap.memory_lock=true
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    ulimits:
      memlock:
        soft: -1
        hard: -1
    volumes:
      - data01:/usr/share/elasticsearch/data
    ports:
      - 9200:9200

  es02:
    image: rkrohk/elasticsearch
    container_name: es02
    environment:
      - node.name=es02
      - cluster.name=es-docker-cluster
      - discovery.seed_hosts=es01  
      - cluster.initial_master_nodes=es01,es02  
      - bootstrap.memory_lock=true
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    ulimits:
      memlock:
        soft: -1
        hard: -1
    volumes:
      - data02:/usr/share/elasticsearch/data


volumes:
  mongodata:
  nltk_data:
  chatbotapi_data:
  data01:
    driver: local

  data02:
    driver: local 
```

Along with the docker compose file, you also need a file named
`.env.gobot` to store secrets.

Structure of `.env.gobot`
```
BOT_TOKEN=<bot token>

BLOCKED_USER=<user ids of banned users>

OWNER=<your telegram user id, get it by using /id>
```
