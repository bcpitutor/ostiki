package repositories

import (
	"github.com/tiki-systems/tikiserver/models"
	"go.uber.org/dig"
)

type BanRepository struct {
	DBLayer models.DBLayer
}

type BanRepositoryResult struct {
	dig.Out
	BanRepository *BanRepository
}

func ProvideBanRepository(db models.DBLayer) BanRepositoryResult {
	return BanRepositoryResult{
		BanRepository: &BanRepository{
			DBLayer: db,
		},
	}
}

func (b *BanRepository) GetBannedUserByEmail(userEmail string) (models.BannedUser, error) {
	return b.DBLayer.GetBannedUserByEmail(userEmail)
}

func (b *BanRepository) GetBannedUsers() ([]models.BannedUser, error) {
	return b.DBLayer.GetBannedUsers()
}

func (b *BanRepository) AddBannedUser(bannedUser models.BannedUser) error {
	return b.DBLayer.AddBannedUser(bannedUser)
}

func (b *BanRepository) UnbanUser(userEmail string) error {
	return b.DBLayer.UnbanUser(userEmail)
}

func (b *BanRepository) IsUserBanned(userEmail string) bool {
	return b.DBLayer.IsUserBanned(userEmail)
}
