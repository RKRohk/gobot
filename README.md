# Gobot

A modular telegram bot written using GoLang. Utilizing microservices architecture


## About

This bot started as a small project to store class notes and keep assignment reminders on telegram. 
Slowly and gradually, a lot of features have been added to it


## Features

### AI

The bot learns from the speaking habits of it's users and learns to talk back. It is surprisingly good, almost seems like a real person is talking back.
The AI feature is handled by a python backend. The bot communicates with it using gRPC. 


### DadJoke

`/dadjoke` - Sends a random Dad joke to the chat

### Echo

Echos back the message sent

### Slap

Slaps the person the command is used against

### Reminder

Keeps reminders

### Save

Saves a note with a tag. Also if a file is given, it is passed through a pipeline where OCR is run on it, 
it is indexed and ready to be searched, like Google! but for your files.

### Get

Gets all notes with the given tag


### Search
`/search #<tag> <query>` - returns the file which contains the notes on query. 
The query need not be exact, as the search engine uses fuzzy algorithms and replaces words with synonyms to find the perfect match. 

