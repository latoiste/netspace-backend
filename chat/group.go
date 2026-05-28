package chat

import "github.com/google/uuid"

type Group struct {
	id        string
	name      string
	hostId    string
	memberIds []string
}

func GenerateGroupId() string {
	id := uuid.New()
	return id.String()
}
