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

