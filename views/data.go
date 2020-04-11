package views

import (
	"lenslocked.com/models"
)

const (
	AlertLvlDanger  = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"
	AlertMsgGeneric = "Something went wrong"
)

type Alert struct {
	Level   string
	Message string
}

type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok {
		d.Alert = &Alert{
			Level:   AlertLvlDanger,
			Message: pErr.Public(),
		}
	} else {
		d.Alert = &Alert{
			Level:   AlertLvlDanger,
			Message: AlertMsgGeneric,
		}
	}
}

type PublicError interface {
	error
	Public() string
}
