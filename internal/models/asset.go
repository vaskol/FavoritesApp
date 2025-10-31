package models

type Asset interface {
	GetID() string
	SetDescription(desc string)
}

// Chart asset
type Chart struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

func (c *Chart) GetID() string              { return c.ID }
func (c *Chart) SetDescription(desc string) { c.Description = desc }

// Insight asset
type Insight struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

func (i *Insight) GetID() string              { return i.ID }
func (i *Insight) SetDescription(desc string) { i.Description = desc }

// Audience asset
type Audience struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Gender      string `json:"gender"`
	Country     string `json:"country"`
	AgeGroup    string `json:"age_group"`
	SocialHours int    `json:"social_hours"`
	Purchases   int    `json:"purchases"`
}

func (a *Audience) GetID() string              { return a.ID }
func (a *Audience) SetDescription(desc string) { a.Description = desc }
