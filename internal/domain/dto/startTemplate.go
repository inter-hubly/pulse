package dto

type StartTemplateDto struct {
	SenderAndReceiver SenderAndReceiverDto `json:"senderAndReceiver"`
	Name              string               `json:"name"`
	Language          string               `json:"language"`
}

type SenderAndReceiverDto struct {
	OwnerNumberId string `json:"OwnerNumberId"`
	From          string `json:"from"`
	To            string `json:"to"`
}
