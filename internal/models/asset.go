package models

// Asset interface
type Asset interface {
	GetID() string
	SetDescription(desc string)
}

// ChartData represents one data point in a chart
type ChartData struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

// Chart asset
type Chart struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	XAxisTitle  string      `json:"x_axis_title"`
	YAxisTitle  string      `json:"y_axis_title"`
	Data        []ChartData `json:"data"`
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
