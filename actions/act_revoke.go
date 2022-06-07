package actions

import (
	"github.com/bcpitutor/ostiki/repositories"
	"go.uber.org/zap"
)

func IsUserRevoked(banRepository *repositories.BanRepository, logger *zap.Logger, userEmail string) bool {
	user, err := banRepository.GetBannedUserByEmail(userEmail)
	if err != nil {
		if err.Error() == "404" {
			return false
		}

		logger.Sugar().Debugf("get banned user list has error for the user: %s", userEmail)
		return true
	}

	if user.UserEmail == userEmail {
		return true
	}

	return false
}
