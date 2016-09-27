package telegrambotgo

import (
	"bytes"
	"encoding/json"
	"net/http"
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

func (b *BotClient) Updates() <-chan Update {
	updatesChannel := make(chan Update)

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
	var message Message
	err := b.executeMethod("sendMessage", request, &message)

	return &message, err
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
