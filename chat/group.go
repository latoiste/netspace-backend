package chat

import "github.com/google/uuid"

type Group struct {
	id        string
	name      string
	hostId    string
	memberIds map[string]bool
}

type GroupInvite struct {
	groupId string
	userId  string
}

func GenerateGroupId() string {
	id := uuid.New()
	return id.String()
}
