package models

// Asset interface
type Asset interface {
	GetID() string
	SetDescription(desc string)
}

// ChartData represents one data point in a chart
type ChartData struct {
	DatapointCode string `json:"datapoint_code"`
	NamespaceCode string `json:"namespace_code"`
	QuestionCode  string `json:"question_code"`
	SuffixCode    string `json:"suffix_code"`
}

// Chart asset
type Chart struct {
	ID          string      `json:"id"`          // Chart unique ID
	Name        string      `json:"name"`        // Chart title
	Description string      `json:"description"` // Chart description
	Attributes  []ChartData `json:"attributes"`  // Essential data points
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
	Name        string `json:"name"`
	Description string `json:"description"` // extra info
}

func (a *Audience) GetID() string              { return a.ID }
func (a *Audience) SetDescription(desc string) { a.Description = desc }
