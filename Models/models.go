package models

import (
	"errors"
	"html"
	"strings"

	"task-5-vix-fullstack/app"
)

type Photo struct {
	ID       int        `gorm:"primary_key;auto_increment" json:"id"`
	Title    string     `gorm:"size:100;not null" json:"title"`
	Caption  string     `gorm:"size:255;not null" json:"caption"`
	PhotoUrl string     `gorm:"size:255;not null;" json:"photo_url"`
	UserId   string     `gorm:"not null" json:"user_id"`
	Author   app.Author `gorm:"author"`
}

// UInisialisasi data User sebelum di save/update
func (p *Photo) Initialize() {
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Caption = html.EscapeString(strings.TrimSpace(p.Caption))
	p.PhotoUrl = html.EscapeString(strings.TrimSpace(p.PhotoUrl))
}

// Validasi data user input sebelum disimpan
func (p *Photo) Validate(action string) error {
	switch strings.ToLower(action) {
	case "upload":
		if p.Title == "" {
			return errors.New("required title")
		} else if p.Caption == "" {
			return errors.New("required caption")
		} else if p.PhotoUrl == "" {
			return errors.New("required photo url")
		} else if p.UserId == "" {
			return errors.New("required user_id")
		}
		return nil
	case "change":
		if p.Title == "" {
			return errors.New("required title")
		} else if p.Caption == "" {
			return errors.New("required caption")
		} else if p.PhotoUrl == "" {
			return errors.New("required photo url")
		}
		return nil
	default:
		return nil
	}
}
