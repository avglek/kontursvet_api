package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"kontursvet-api/internal/model"
)

type MaxBotService struct {
	Token  string
	ChatID string
	Active bool
}

func NewMaxBotService(token, chatID string, active bool) *MaxBotService {
	return &MaxBotService{
		Token:  token,
		ChatID: chatID,
		Active: active,
	}
}

func (s *MaxBotService) SendLeadNotification(lead *model.Lead) error {
	if !s.Active || s.Token == "" || s.ChatID == "" {
		return nil
	}

	message := fmt.Sprintf(
		"📋 Новый лид!\n\n"+
			"👤 Имя: %s\n"+
			"📞 Телефон: %s\n"+
			"🏠 Тип: %s\n"+
			"📍 Локация: %s\n"+
			"💬 Сообщение: %s",
		lead.Name,
		lead.PhoneFormat,
		lead.HomeType,
		lead.Location,
		lead.Message,
	)

	payload := map[string]interface{}{
		"chat_id":    s.ChatID,
		"text":       message,
		"parse_mode": "HTML",
	}

	return s.sendMessage(payload)
}

func (s *MaxBotService) sendMessage(payload map[string]interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("max bot marshal: %w", err)
	}

	req, err := http.NewRequest("POST",
		"https://max.ru/api/v1/bot/"+s.Token+"/send_message",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("max bot request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("max bot send: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("max bot error status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
