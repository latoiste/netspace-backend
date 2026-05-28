package auth

type Blacklist struct {
	blacklist map[string]bool
}

func NewBlacklist() *Blacklist {
	return &Blacklist{
		blacklist: make(map[string]bool),
	}
}

func (b *Blacklist) Add(token string) {
	b.blacklist[token] = true
}

func (b *Blacklist) IsBlacklisted(token string) bool {
	_, ok := b.blacklist[token]
	return ok
}
