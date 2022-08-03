package models

type MailEvent struct {
	TaskID      int32  `json:"task_id"`
	Description string `json:"description"`
	Body        string `json:"body"`
	Addressee   string `json:"addressee"`
	MailType    string `json:"mail_type"`
	ApproveLink string `json:"approve_link"`
	RejectLink  string `json:"reject_link"`
}
