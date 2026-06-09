package api

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/latoiste/netspace/model"
)

type UserDTO struct {
	Id         string        `json:"id"`
	Slug       string        `json:"slug"`
	Name       string        `json:"name"`
	Emoji      string        `json:"emoji"`
	Occupation string        `json:"occupation"`
	Interests  []InterestDTO `json:"interests"`
}

type InterestDTO struct {
	Emoji string `json:"emoji"`
	Label string `json:"label"`
}

type InterestPercentageDTO struct {
	InterestDTO
	Percentage int `json:"percentage"`
}

type ActiveUserDTO struct {
	Id        string        `json:"id"`
	Name      string        `json:"name"`
	Avatar    string        `json:"avatar"`
	Gender    string        `json:"gender"`
	Interests []InterestDTO `json:"interests"`
	Duration  int           `json:"duration"`
}

type AdminDTO struct {
	Name   string `json:"name"`
	Role   string `json:"role"`
	Plan   string `json:"plan"`
	Avatar string `json:"avatar"`
}

type NotificationDTO struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	Emoji           string `json:"emoji"`
	AvatarGradient  string `json:"avatarGradient"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	TimestampString string `json:"timestamp"`
	Unread          bool   `json:"unread"`
	PrimaryLabel    string `json:"primaryLabel,omitempty"`
	SecondaryLabel  string `json:"secondaryLabel,omitempty"`
	GroupId         string `json:"groupId,omitempty"`
	SenderId        string `json:"senderId,omitempty"`
}

type AnalyticsDTO struct {
	Label     string `json:"label"`
	Value     string `json:"value"`
	Delta     string `json:"delta"`
	DeltaType string `json:"deltaType"`
}

type HourlyCheckInDTO struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

type MessageDTO struct {
	Id             string `json:"id"`
	Type           string `json:"type"`
	Name           string `json:"name"`
	Emoji          string `json:"emoji"`
	AvatarGradient string `json:"avatarGradient"`
	LastMessage    string `json:"lastMessage"`
	Timestamp      string `json:"timestamp"`
	Unread         bool   `json:"unread"`
	Subtitle       string `json:"subtitle,omitempty"`
}

func ConstructUserDTO(user model.User) UserDTO {
	interestDTOs := make([]InterestDTO, 0, len(user.Interests))

	for _, interest := range user.Interests {
		interestDTO := ConstructInterestDTO(interest)
		interestDTOs = append(interestDTOs, interestDTO)
	}

	return UserDTO{
		Id:         user.Id,
		Slug:       user.Slug,
		Name:       user.Name,
		Emoji:      EmojiForUser(user.Id),
		Occupation: user.Occupation,
		Interests:  interestDTOs,
	}
}

func ConstructAdminDTO(admin model.Admin) AdminDTO {
	return AdminDTO{
		Name:   admin.Name,
		Role:   admin.Role,
		Plan:   admin.Plan,
		Avatar: admin.Avatar,
	}
}

func ConstructInterestDTO(interest model.Interest) InterestDTO {
	return InterestDTO{
		Emoji: interest.Emoji,
		Label: interest.Label,
	}
}

func ConstructActiveUserDTO(user model.User) ActiveUserDTO {
	interestDTOs := make([]InterestDTO, 0, len(user.Interests))

	for _, interest := range user.Interests {
		interestDTO := ConstructInterestDTO(interest)
		interestDTOs = append(interestDTOs, interestDTO)
	}

	duration := time.Since(user.CreatedAt).Minutes()

	return ActiveUserDTO{
		Id:        user.Id,
		Name:      user.Name,
		Avatar:    EmojiForUser(user.Id),
		Gender:    user.Gender,
		Interests: interestDTOs,
		Duration:  int(math.Floor(duration)),
	}
}

func ConstructNotificationDTO(notif model.Notification) NotificationDTO {
	return NotificationDTO{
		Id:              notif.Id,
		Type:            notif.Type,
		Emoji:           notif.Emoji,
		AvatarGradient:  notif.AvatarGradient,
		Title:           notif.Title,
		Description:     notif.Description,
		TimestampString: notif.Timestamp.Local().Format("15:04"),
		Unread:          notif.Unread,
		PrimaryLabel:    notif.PrimaryLabel,
		SecondaryLabel:  notif.SecondaryLabel,
		GroupId:         notif.GroupId,
		SenderId:        notif.SenderId,
	}
}

func ConstructAnalyticsDTO(prevValue int, currentValue int, label string, deltaType string) AnalyticsDTO {
	var deltaPercentage int
	if total := prevValue + currentValue; total != 0 {
		deltaPercentage = int(math.Round(float64(currentValue) / float64(total) * 100))
	}

	var delta string

	switch deltaType {
	case "down":
		delta = fmt.Sprintf("↓ %v%% vs kemarin", deltaPercentage)
	case "up":
		delta = fmt.Sprintf("↑ %v%% vs kemarin", deltaPercentage)
	case "live":
		delta = "Live"
	}

	return AnalyticsDTO{
		Label:     label,
		Value:     fmt.Sprint(currentValue),
		Delta:     delta,
		DeltaType: deltaType,
	}
}

func ConstructPrivateMessageDTO(msg model.PrivateMessage, recipient model.User) MessageDTO {
	return MessageDTO{
		Id:             msg.RecipientId,
		Type:           "dm",
		Name:           recipient.Name,
		Emoji:          EmojiForUser(recipient.Id),
		AvatarGradient: "linear-gradient(135deg, rgba(56,100,255,0.5), rgba(100,60,255,0.4))",
		LastMessage:    msg.Message,
		Timestamp:      msg.Timestamp.Local().Format("15:04"),
	}
}

// ── Chat history DTOs (REST: load past messages on chat open) ──

type ChatMessageDTO struct {
	Id        string `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	IsMine    bool   `json:"isMine"`
	IsRead    bool   `json:"isRead"`
}

// PublicChatMessageDTO mirrors the live new_public_message wire shape so the
// public room can render REST history and live events with the same code path.
type PublicChatMessageDTO struct {
	Id          string `json:"id"`
	SenderId    string `json:"senderId"`
	SenderName  string `json:"senderName"`
	SenderEmoji string `json:"senderEmoji"`
	Message     string `json:"message"`
	Timestamp   string `json:"timestamp"`
	IsMine      bool   `json:"isMine"`
}

type GroupChatMessageDTO struct {
	Id          string `json:"id"`
	SenderId    string `json:"senderId"`
	SenderName  string `json:"senderName"`
	SenderEmoji string `json:"senderEmoji"`
	Message     string `json:"message"`
	Timestamp   string `json:"timestamp"`
	IsMine      bool   `json:"isMine"`
}

type DMPartnerDTO struct {
	Name       string `json:"name"`
	Emoji      string `json:"emoji"`
	Occupation string `json:"occupation"`
	Interests  string `json:"interests"`
	IsOnline   bool   `json:"isOnline"`
}

type GroupMemberDTO struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Emoji  string `json:"emoji"`
	IsHost bool   `json:"isHost"`
}

type GroupSummaryDTO struct {
	Name  string `json:"name"`
	Emoji string `json:"emoji"`
}

// InterestsToString renders interests as "☕ Kopi, 📚 Buku" for the DM header.
func InterestsToString(interests []model.Interest) string {
	parts := make([]string, 0, len(interests))
	for _, i := range interests {
		parts = append(parts, fmt.Sprintf("%s %s", i.Emoji, i.Label))
	}
	return strings.Join(parts, ", ")
}

func ConstructChatMessageDTO(msg model.PrivateMessage, myId string) ChatMessageDTO {
	return ChatMessageDTO{
		Id:        msg.MessageId,
		Message:   msg.Message,
		Timestamp: msg.Timestamp.Local().Format("15:04"),
		IsMine:    msg.SenderId == myId,
		IsRead:    msg.IsRead,
	}
}

func ConstructDMPartnerDTO(partner model.User) DMPartnerDTO {
	return DMPartnerDTO{
		Name:       partner.Name,
		Emoji:      EmojiForUser(partner.Id),
		Occupation: partner.Occupation,
		Interests:  InterestsToString(partner.Interests),
		IsOnline:   partner.IsActive,
	}
}

func ConstructGroupMessageDTO(msg model.GroupMessage, group model.Group) MessageDTO {
	return MessageDTO{
		Id:             group.Id,
		Type:           "group",
		Name:           group.Name,
		Emoji:          "☕",
		AvatarGradient: "linear-gradient(135deg, rgba(56,100,255,0.5), rgba(100,60,255,0.4))",
		LastMessage:    msg.Message,
		Timestamp:      msg.Timestamp.Local().Format("15:04"),
		Subtitle:       "Group session",
	}

}

// ConstructGroupSummaryDTO builds a chat-list entry for a group the user belongs
// to, whether or not it has any messages yet. A group with no messages shows a
// gentle placeholder so it still appears in the list.
func ConstructGroupSummaryDTO(groupId, name, lastMessage string, hasMessage bool, ts time.Time, hasTimestamp bool) MessageDTO {
	preview := "Belum ada pesan"
	if hasMessage {
		preview = lastMessage
	}
	timestamp := ""
	if hasTimestamp {
		timestamp = ts.Local().Format("15:04")
	}

	return MessageDTO{
		Id:             groupId,
		Type:           "group",
		Name:           name,
		Emoji:          "☕",
		AvatarGradient: "linear-gradient(135deg, rgba(56,100,255,0.5), rgba(100,60,255,0.4))",
		LastMessage:    preview,
		Timestamp:      timestamp,
		Subtitle:       "Group session",
	}
}
