package types

import "time"

type Setting struct {
	Key   string `gorm:"column:key;index;unique;size:255" json:"key"`
	Value string `gorm:"column:value;null;type:text;size:20000" json:"value"`
	//Title       string    `gorm:"column:title;null" json:"title"`
	Description string    `gorm:"column:description;null;size:500" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (p *Setting) TableName() string {
	return "setting"
}
