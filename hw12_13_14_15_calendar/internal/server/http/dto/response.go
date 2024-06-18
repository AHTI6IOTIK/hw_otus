package dto

import "github.com/AHTI6IOTIK/hw_otus/hw12_13_14_15_calendar/internal/storage"

type CreateReply struct {
	ID string `json:"id"`
}

type ListReply struct {
	Result []storage.Event `json:"result"`
}

type Result struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
