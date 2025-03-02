package models

type KafkaEvent struct {
	EventType string `json:"event_type"` // "create", "update", "delete"
	Book      Book   `json:"book"`
}
type Book struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
	Year   int    `json:"year" binding:"required"`
}
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Object  interface{} `json:"object,omitempty"`
}

func (Book) TableName() string {
	return "books"
}
