package repository

import "github.com/gofrs/uuid"

// StarRepository チャンネルスターリポジトリ
type StarRepository interface {
	AddStar(userID, channelID uuid.UUID) error
	RemoveStar(userID, channelID uuid.UUID) error
	GetStaredChannels(userID uuid.UUID) ([]uuid.UUID, error)
}
