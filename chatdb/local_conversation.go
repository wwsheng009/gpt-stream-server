package chatdb

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
)

type LocalConversation struct {
	Setting *ApiSetting
}

const cacheLocation = "./data/cache"

func (l *LocalConversation) LoadApiSetting() (*ApiSetting, error) {

	data, err := os.ReadFile("./gpt.config.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		panic(err)
		// return nil, err
	}
	setting := ApiSetting{}
	err = json.Unmarshal(data, &setting)
	if err != nil {
		return nil, err
	}
	if setting.MaxSendLines > 20 {
		setting.MaxSendLines = 20
	}
	if setting.MaxSendLines < 1 {
		setting.MaxSendLines = 1
	}
	l.Setting = &setting
	return &setting, nil
}

func (l *LocalConversation) Save(uuid string, c *Conversation) error {

	fname := path.Join(cacheLocation, uuid) + ".json"
	// Convert the struct to JSON
	jsonData, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Write the JSON data to a file
	file, err := os.Create(fname)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (l *LocalConversation) CreateNewconversation(title string) (*Conversation, error) {
	err := makeFolder()
	if err != nil {
		return nil, err
	}
	id := uuid.New().String()

	c := Conversation{}
	c.ConversationId = id
	c.Title = title
	l.Save(id, &c)
	return &c, nil
}

func (l *LocalConversation) FindConversationById(conversationId string) (*Conversation, error) {
	fname := path.Join(cacheLocation, conversationId) + ".json"

	var p Conversation
	// Read the JSON data from the file
	jsonData, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Convert the JSON data to a struct

	err = json.Unmarshal(jsonData, &p)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	p.Messages = GetLastLines(p.Messages, l.Setting.MaxSendLines)

	return &p, nil
}

func (l *LocalConversation) CreateNewMessage(conversation *Conversation, prompt string, answer string, seconds float64) error {
	fname := path.Join(cacheLocation, conversation.ConversationId) + ".json"

	var c Conversation
	// Read the JSON data from the file
	jsonData, err := os.ReadFile(fname)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, &c)
	if err != nil {
		return err
	}

	mes := ConvMessage{Prompt: prompt, Completion: answer, Seconds: seconds}
	c.Messages = append(c.Messages, mes)
	err = l.Save(c.ConversationId, &c)
	if err != nil {
		return err
	}
	return nil
}

func makeFolder() error {

	_, err := os.Stat(cacheLocation)
	if os.IsNotExist(err) {
		// Create the folder if it does not exist
		err := os.MkdirAll(cacheLocation, 0755)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("Folder created successfully")
	} else {
		fmt.Println("Folder already exists")
		return nil
	}
	return nil
}
