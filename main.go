package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	botAccessToken := "2056045746:AAEHVepiBuHuBHTSmN-kBlGDaSDCBbEMWmk"
	botApi := "https://api.telegram.org/bot"
	botRequest := botApi + botAccessToken
	offset := 0
	for {
		updates, err := getUpdates(botRequest, offset)
		if err != nil {
			log.Println("Error", err.Error())
		}
		fmt.Println(updates)
		for _, update := range updates {
			respond(botRequest, update)
			offset = update.UpdateId + 1
		}
	}
}

// Запрос обновлений
func getUpdates(botRequest string, offset int) ([]Update, error) {
	resp, err := http.Get(botRequest + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response RestResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return response.Result, nil
}

// Получить файл с голосовым сообщением
func getFile(botRequest string, message BotMessage) error {

	resp, err := http.Get(botRequest + "/getFile?file_id=" + message.Voice.FileId)

	println("ID Файла: " + message.Voice.FileId)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var fileResponse FileResponse
	err = json.Unmarshal(body, &fileResponse)

	if err != nil {
		return err
	}

	println("Путь файла: " + fileResponse.Result.Path)

	resp, err = http.Get("https://api.telegram.org/file/bot" + "2056045746:AAEHVepiBuHuBHTSmN-kBlGDaSDCBbEMWmk" + "/" + fileResponse.Result.Path)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create("test.oga")

	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return err
	}

	return nil
}

// Ответ на обновления
func respond(botRequest string, update Update) error {
	var message = BotMessage{
		ChatId: update.Message.Chat.ChatId,
		Text:   update.Message.Text,
		Voice:  update.Message.Voice,
	}

	if message.Voice.Duration != 0 {
		message.Text = "Голосовое сообщение"
		getFile(botRequest, message)
	}

	buffer, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = http.Post(botRequest+"/sendMessage", "application/json", bytes.NewBuffer(buffer))

	if err != nil {
		return err
	}

	return nil
}
