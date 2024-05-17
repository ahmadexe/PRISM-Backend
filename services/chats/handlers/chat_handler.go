package handlers

import "github.com/ahmadexe/prism-backend/services/chats/repository"


type ChatHandler struct {
	repo *repository.ChatRepo
}

func InitChatHandler(repo *repository.ChatRepo) *ChatHandler {
	return &ChatHandler{repo: repo}
}