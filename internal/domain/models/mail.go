package models

type MailEvent struct {
	TaskId      int32  `json:"task_id"`
	Description string `json:"description"`
	Addressee   string `json:"addressee"`
	MailType    string `json:"mail_type"`
	ApproveLink string `json:"approve_link"`
	RejectLink  string `json:"reject_link"`
}
