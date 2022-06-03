package repositories

import (
	"github.com/tiki-systems/tikiserver/models"
	"go.uber.org/dig"
)

type GroupRepository struct {
	DBLayer    models.DBLayer
	CacheLayer models.CacheLayer
}

type GroupRepositoryResult struct {
	dig.Out
	GroupRepository *GroupRepository
}

func ProvideGroupRepository(db models.DBLayer, cl models.CacheLayer) GroupRepositoryResult {
	return GroupRepositoryResult{
		GroupRepository: &GroupRepository{
			DBLayer:    db,
			CacheLayer: cl,
		},
	}
}

func (g *GroupRepository) IsUserInTikiadmins(userEmail string) bool {
	return g.DBLayer.IsUserInTikiadmins(userEmail)
}

func (g *GroupRepository) GetAllGroups() ([]models.TicketGroup, error) {
	return g.DBLayer.GetAllGroups()
}

func (g *GroupRepository) GetGroup(groupName string) (models.TicketGroup, error) {
	return g.DBLayer.GetGroup(groupName)
}

func (g *GroupRepository) DeleteGroup(groupName string) error {
	return g.DBLayer.DeleteGroup(groupName)
}

func (g *GroupRepository) DoesGroupExist(groupName string) bool {
	return g.DBLayer.DoesGroupExist(groupName)
}

func (g *GroupRepository) IsUserMemberOfGroup(member string, groupName string) bool {
	return g.DBLayer.IsUserMemberOfGroup(member, groupName)
}

func (g *GroupRepository) CreateGroup(group models.TicketGroup) error {
	return g.DBLayer.CreateGroup(group)
}

func (g *GroupRepository) AddMemberToGroup(newMember string, groupName string, changedBy string) error {
	return g.DBLayer.AddMemberToGroup(newMember, groupName, changedBy)
}

func (g *GroupRepository) DelMemberFromGroup(memberToDelete string, groupName string, changedBy string) error {
	return g.DBLayer.DelMemberFromGroup(memberToDelete, groupName, changedBy)
}
