package telegrambotgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

const baseBotURI string = "https://api.telegram.org/bot"

type ClientConfig struct {
	Token              string
	UpdatesChannelSize int
	ErrorsChannelSize  int
}

type BotClient struct {
	config         *ClientConfig
	closed         bool
	updatesOffset  int
	updatesChannel chan Update
	errorsChannel  chan error
}

func NewClientConfig(token string) *ClientConfig {
	return &ClientConfig{token,
		100,
		100,
	}
}

func NewBotClient(config *ClientConfig) *BotClient {
	return &BotClient{config,
		false,
		0,
		make(chan Update, config.UpdatesChannelSize),
		make(chan error, config.ErrorsChannelSize),
	}
}

func (b *BotClient) Close() {
	b.closed = true
}

func (b *BotClient) Errors() <-chan error {
	return b.errorsChannel
}

func (b *BotClient) Updates() <-chan *Update {
	updatesChannel := make(chan *Update)

	go func() {
		// Get updates received by this bot
		for !b.closed {
			var response UpdatesResponse

			// Build the parameters for this request
			parameters := make(map[string]interface{})
			parameters["offset"] = b.updatesOffset
			parameters["timeout"] = 10

			err := b.executeMethod("getUpdates", parameters, &response)

			if err == nil && response.Ok {
				for _, update := range response.Result {
					b.updatesOffset = update.UpdateId + 1
					updatesChannel <- update
				}
			} else {
				b.notifyError(err)
			}
		}
	}()

	return updatesChannel
}

func (b *BotClient) SendMessage(request *SendMessageRequest) (*Message, error) {
	var messageResponse SendMessageResponse
	err := b.executeMethod("sendMessage", request, &messageResponse)

	if err != nil {
		log.Printf("Error processing response %v\n", err)
		return nil, err
	}

	if !messageResponse.Ok {
		return nil, errors.New("Error sending message")
	}

	return messageResponse.Result, nil
}

func (b *BotClient) SendPhoto(request *SendPhotoRequest) (*Message, error) {
	var messageResponse SendMessageResponse

	if request.IsLocalFile {
		// Send a multipart request with the specified field
		parameters := make(map[string]string)
		parameters["chat_id"] = strconv.FormatInt(request.ChatId, 10)

		// TODO Add other possible parameters

		err := b.executeMethodMultipart("sendPhoto", request.Photo, "photo", parameters, &messageResponse)
		if err != nil {
			log.Printf("Error processing response %v\n", err)
			return nil, err
		}

	} else {
		// The file is in the server, send a regular json request
		err := b.executeMethod("sendPhoto", request, &messageResponse)
		if err != nil {
			log.Printf("Error processing response %v\n", err)
			return nil, err
		}
	}

	return messageResponse.Result, nil
}

func (b *BotClient) notifyError(err error) {
	if err != nil && len(b.errorsChannel) < b.config.ErrorsChannelSize {
		b.errorsChannel <- err
	}
}

func (b *BotClient) executeMethod(method string, parameters interface{}, result interface{}) error {

	var jsonParams, err = json.Marshal(parameters)
	if err != nil {
		return err
	}

	// Connect to bot api to retrieve the updates
	resp, err := http.Post(baseBotURI+b.config.Token+"/"+method, "application/json", bytes.NewBuffer(jsonParams))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Decode response
	decoder := json.NewDecoder(resp.Body)

	return decoder.Decode(result)
}

func (b *BotClient) executeMethodMultipart(method string, file_path string, multipart_file_field string, extra_params map[string]string, result interface{}) (err error) {
	var requestBytes bytes.Buffer
	w := multipart.NewWriter(&requestBytes)

	// Open the file
	f, err := os.Open(file_path)
	if err != nil {
		return
	}
	defer f.Close()

	// Create the file field
	fw, err := w.CreateFormFile(multipart_file_field, file_path)
	if err != nil {
		return
	}

	// Copy the contents into the multipart file field
	if _, err = io.Copy(fw, f); err != nil {
		return
	}

	// Now create the other fields
	for k, v := range extra_params {

		// In case there's a key equal to the multipart file key, just ignore it and continue
		if k == multipart_file_field {
			continue
		}

		if fw, err = w.CreateFormField(k); err != nil {
			return
		}

		if _, err = fw.Write([]byte(v)); err != nil {
			return
		}
	}

	// This is necessary in order to set the terminating boundary
	w.Close()

	req, err := http.NewRequest("POST", baseBotURI+b.config.Token+"/"+method, &requestBytes)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode == http.StatusBadRequest {
		err = errors.New("Bad request sent to server")
		return
	}

	defer resp.Body.Close()

	// Decode response
	decoder := json.NewDecoder(resp.Body)

	return decoder.Decode(result)

	return
}
