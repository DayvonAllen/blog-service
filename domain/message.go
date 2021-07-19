package domain

// Message messageType 201 user created
// messageType 200 user updated
type Message struct {
	Post        Post   `form:"Post" json:"Post"`
	Tag        Tag   `form:"Tag" json:"Tag"`
	Event        Event  `form:"Event" json:"Event"`
	MessageType  int    `form:"messageType" json:"messageType"`
	ResourceType string `form:"resourceType" json:"resourceType"`
}
