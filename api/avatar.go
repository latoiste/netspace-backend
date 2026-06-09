package api

// avatarEmojis is a small palette of friendly avatar faces. Users don't pick
// their own avatar during check-in, so we assign one deterministically from
// their (stable) userId. This keeps avatars consistent across REST + WebSocket
// without needing an extra DB column.
var avatarEmojis = []string{
	"🦊", "🐼", "🐧", "🦉", "🐙", "🦁", "🐯", "🐸",
	"🐵", "🦄", "🐳", "🐝", "🦋", "🐺", "🐶", "🐱",
	"🐰", "🐻", "🐨", "🦝", "🐢", "🦓", "🐬", "🦜",
}

// EmojiForUser returns a stable avatar emoji for a given user id.
func EmojiForUser(userId string) string {
	if userId == "" {
		return "🙂"
	}

	// Simple deterministic hash over the id bytes.
	var sum uint32
	for i := 0; i < len(userId); i++ {
		sum = sum*31 + uint32(userId[i])
	}

	return avatarEmojis[sum%uint32(len(avatarEmojis))]
}
