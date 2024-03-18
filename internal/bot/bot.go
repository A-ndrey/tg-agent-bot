package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type Bot struct {
	url    string
	chatID int
	client *http.Client
}

func New(token string, chatID int) *Bot {
	return &Bot{
		url:    fmt.Sprintf("https://api.telegram.org/bot%s", token),
		chatID: chatID,
		client: &http.Client{Timeout: time.Minute},
	}
}

func (b *Bot) GetChatID() (int, error) {
	resp, err := b.client.Get(b.url + "/getUpdates")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("%s: %w", resp.Status, parseErr(resp.Body))
	}

	respStruct := struct {
		Result []struct {
			Message struct {
				Chat struct {
					ID int `json:"id"`
				} `json:"chat"`
			} `json:"message"`
		} `json:"result"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		return 0, err
	}

	return respStruct.Result[len(respStruct.Result)-1].Message.Chat.ID, err
}

func (b *Bot) SendMessage(msg string) error {
	req := struct {
		ChatID int    `json:"chat_id"`
		Text   string `json:"text"`
	}{
		ChatID: b.chatID,
		Text:   msg,
	}

	encReq, err := json.Marshal(req)
	if err != nil {
		return err
	}

	response, err := b.client.Post(b.url+"/sendMessage", "application/json", bytes.NewReader(encReq))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %w", response.Status, parseErr(response.Body))
	}

	return nil
}

func (b *Bot) SendFile(reader io.Reader, fileName string) error {
	body := bytes.Buffer{}
	w := multipart.NewWriter(&body)
	if err := w.WriteField("chat_id", strconv.Itoa(b.chatID)); err != nil {
		return err
	}
	part, err := w.CreateFormFile("document", fileName)
	if err != nil {
		return err
	}

	if _, err := io.Copy(part, reader); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	response, err := b.client.Post(b.url+"/sendDocument", w.FormDataContentType(), &body)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %w", response.Status, parseErr(response.Body))
	}

	return nil
}

func parseErr(r io.Reader) error {
	respErr := struct {
		Description string `json:"description"`
	}{}

	if err := json.NewDecoder(r).Decode(&respErr); err != nil {
		return err
	}

	return errors.New(respErr.Description)
}
