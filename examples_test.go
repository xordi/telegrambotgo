package telegrambotgo_test

import (
	"fmt"

	"github.com/xordi/telegrambotgo"
)

func ExampleSendMessage() {
	config := telegrambotgo.NewClientConfig("72918624:AAG4Sj4k6zDG6-g7YSwZqJ5Psmx4b1hGrZM")
	client := telegrambotgo.NewBotClient(config)

	request := new(telegrambotgo.SendMessageRequest)
	request.ChatId = 4248038
	request.Text = "This is a test"

	msg, err := client.SendMessage(request)

	if err == nil {
		fmt.Println(msg.Text)
	} else {
		fmt.Printf("Error: %v", err)
	}

	// Output: This is a test
}

func ExampleSendPhoto() {
	config := telegrambotgo.NewClientConfig("72918624:AAG4Sj4k6zDG6-g7YSwZqJ5Psmx4b1hGrZM")
	client := telegrambotgo.NewBotClient(config)

	request := new(telegrambotgo.SendPhotoRequest)
	request.ChatId = 4248038
	request.IsLocalFile = true
	request.Photo = "examples/test.jpg"
	request.Caption = "This is a photo"

	msg, err := client.SendPhoto(request)

	if err == nil {
		fmt.Printf("Photo sent to chat %d\n", msg.Chat.ChatId)
	} else {
		fmt.Printf("Error: %v", err)
	}

	// Output: Photo sent to chat 4248038
}
