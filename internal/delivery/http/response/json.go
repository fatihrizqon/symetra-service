package response

import "github.com/google/uuid"

type AuthJSON struct {
	Status      int      `json:"status"`
	Message     string   `json:"message"`
	User        UserInfo `json:"user"`
	AccessToken string   `json:"access_token"`
}

type UserInfo struct {
	Id              uuid.UUID `json:"id"`
	Username        string    `json:"username"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Status          int       `json:"status"`
	IsSuperadmin    bool      `json:"is_superadmin"`
	EmailVerifiedAt string    `json:"email_verified_at"`
	Roles           []string  `json:"roles"`
	Permissions     []string  `json:"permissions"`
}

type SelectJSON struct {
	Data interface{} `json:"data,omitempty"`
}

type JSON struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type Meta struct {
	Search     string `json:"search,omitempty"`
	Info       string `json:"info"`
	Page       int    `json:"page"`
	TotalCount int    `json:"total_count"`
	TotalPages int    `json:"total_pages"`
	PageSize   int    `json:"page_size"`
	Links      Links  `json:"links"`
}

type Links struct {
	CurrentPage string  `json:"current_page"`
	FirstPage   string  `json:"first_page"`
	LastPage    string  `json:"last_page"`
	NextPage    *string `json:"next_page,omitempty"`
	PrevPage    *string `json:"prev_page,omitempty"`
}
