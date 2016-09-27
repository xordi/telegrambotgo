package telegrambotgo

type ApiResponse struct {
	Ok bool `json:"ok"`
}

type UpdatesResponse struct {
	ApiResponse
	Result []*Update `json:"result"`
}

type SendMessageResponse struct {
	ApiResponse
	Result *Message `json:"result"`
}

type Update struct {
	UpdateId      int      `json:"update_id"`
	Message       *Message `json:"message"`
	EditedMessage *Message `json:"edited_message"`
}

type Message struct {
	MessageId       int              `json:"message_id"`
	From            *User            `json:"from"`
	Date            int              `json:"date"`
	Chat            *Chat            `json:"chat"`
	ForwardFrom     *User            `json:"forward_from"`
	ForwardFromChat *Chat            `json:"forward_from_chat"`
	Text            string           `json:"text"`
	Entities        []*MessageEntity `json:"entities"`
	NewChatMember   *User            `json:"new_chat_member"`
	LeftChatMember  *User            `json:"left_chat_member"`
	Caption         string           `json:"caption"`
	Photo           []*PhotoSize     `json:"photo"`
}

type User struct {
	UserId    int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
}

type Chat struct {
	ChatId    int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	UserName  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Url    string `json:"url"`
	User   *User  `json:"user"`
}

type PhotoSize struct {
	FileId   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size"`
}

type SendMessageRequest struct {
	ChatId                int64  `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview,omitempty"`
	DisableNotification   bool   `json:"disable_notification,omitempty"`
	ReplyToMessageId      int    `json:"reply_to_message_id,omitempty"`
}

// TODO Add other possible parameters
type SendPhotoRequest struct {
	ChatId      int64  `json:"chat_id"`
	Photo       string `json:"photo"` // This might be a byte representation of the photo or a string containing the file id
	IsLocalFile bool   `json:"-"`
}
