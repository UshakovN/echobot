package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Модель бота
type Model struct {
	botAccessToken string
	botApi         string
	botRequest     string
	yandexApiKey   string
	yandexRequest  string
}

// Точка старта бота
func main() {
	var model Model
	model.botAccessToken = "bot2056045746:AAEHVepiBuHuBHTSmN-kBlGDaSDCBbEMWmk"
	model.botApi = "https://api.telegram.org/"
	model.botRequest = model.botApi + model.botAccessToken
	model.yandexApiKey = "AQVNwbhyxMeFNrOcHO6sthjsOgd5gzq3EPhIkCir"
	model.yandexRequest = "https://stt.api.cloud.yandex.net/speech/v1/stt:recognize?"

	offset := 0
	for {
		updates, err := getUpdates(model, offset)
		if err != nil {
			log.Println("Error", err.Error())
		}
		fmt.Println(updates)
		for _, update := range updates {
			respond(model, update)
			offset = update.UpdateId + 1
		}
	}
}

// Запрос обновлений
func getUpdates(model Model, offset int) ([]Update, error) {
	resp, err := http.Get(model.botRequest + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response UpdateResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return response.Result, nil
}

// Получить файл с голосовым сообщением
func getFile(model Model, message BotMessage) (string, error) {

	var transcription string = "uncorrect"
	respPath, err := http.Get(model.botRequest + "/getFile?file_id=" + message.Voice.FileId)
	if err != nil {
		return transcription, err
	}
	println("ID Файла: " + message.Voice.FileId)

	defer respPath.Body.Close()
	bodyPath, err := ioutil.ReadAll(respPath.Body)
	if err != nil {
		return transcription, err
	}

	var fileResponse FileResponse
	err = json.Unmarshal(bodyPath, &fileResponse)
	if err != nil {
		return transcription, err
	}

	println("Путь файла: " + fileResponse.Result.Path)
	respFile, err := http.Get(model.botApi + "/file/" + model.botAccessToken + "/" + fileResponse.Result.Path)
	if err != nil {
		return transcription, err
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	param := strings.Join([]string{"topic=general", "lang=ru-RU"}, "&")

	println("Параметры запроса: " + param)
	println("Пример запроса: " + model.yandexRequest + param)
	println("Апи-Ключ Яндекса: " + model.yandexApiKey)

	req, err := http.NewRequest(http.MethodPost, model.yandexRequest+param, respFile.Body)
	if err != nil {
		return transcription, err
	}

	req.Header.Add("Authorization", "Api-Key "+model.yandexApiKey)
	respTranscription, err := client.Do(req)
	if err != nil {
		return transcription, err
	}

	defer respTranscription.Body.Close()
	bodyTransription, err := ioutil.ReadAll(respTranscription.Body)
	if err != nil {
		return transcription, err
	}

	var yandexResponse YandexResponse
	err = json.Unmarshal(bodyTransription, &yandexResponse)
	if err != nil {
		return transcription, err
	}

	transcription = yandexResponse.Result

	/*

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

	*/

	return transcription, nil
}

// Ответ на обновления
func respond(model Model, update Update) error {
	var message = BotMessage{
		ChatId: update.Message.Chat.ChatId,
		Text:   update.Message.Text,
		Voice:  update.Message.Voice,
	}

	if message.Voice.Duration != 0 {
		output, err := getFile(model, message)
		if err != nil {
			return err
		}
		message.Text = output

		// message.Text = "Голосовое сообщение"

	}

	buffer, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = http.Post(model.botRequest+"/sendMessage", "application/json", bytes.NewBuffer(buffer))

	if err != nil {
		return err
	}

	return nil
}
