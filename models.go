package main

// Модели Telegram API

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat  Chat   `json:"chat"`
	Text  string `json:"text"`
	Voice Voice  `json:"voice"`
}

type Voice struct {
	FileId   string `json:"file_id"`
	Duration int    `json:"duration"`
}

type FileResponse struct {
	Result File `json:"result"`
}

type File struct {
	Path string `json:"file_path"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type UpdateResponse struct {
	Result []Update `json:"result"`
}

type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
	Voice  Voice  `json:"voice"`
}

type YandexResponse struct {
	Result string `json:"result"`
}
