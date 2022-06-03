package actions

// func DoesUserToAccessTicket(groupRepository *repositories.GroupRepository, userEmail, ticketPath string) bool {
// 	if isUserInTikiadmins(groupRepository, userEmail) {
// 		return true
// 	}

// 	uGroupMembership, err := getUserGroupMembershipNames(userEmail)
// 	if err != nil {
// 		return false
// 	}

// 	tGroupOwners, err := getTicketGroupOwners(ticketPath)
// 	if err != nil {
// 		return false
// 	}

// 	for _, gOwner := range tGroupOwners {
// 		for _, membership := range uGroupMembership {
// 			if gOwner == membership {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

// func getUserGroupMembershipNames(groupRepository *repositories.GroupRepository, userEmail string) ([]string, error) {
// 	var userGroupNamesList []string
// 	// search groups and group members

// 	existingGroups, err := groupRepository.GetAllGroups()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, group := range existingGroups {
// 		gMembers := group.GroupMembers
// 		for _, member := range gMembers {
// 			if userEmail == member {
// 				userGroupNamesList = append(userGroupNamesList, group.GroupName)
// 			}
// 		}
// 	}
// 	return userGroupNamesList, nil
// }

// func isUserInTikiadmins(groupRepository *repositories.GroupRepository, userEmail string) bool {
// 	uGroupList, err := getUserGroupMembershipNames(groupRepository, userEmail)
// 	if err != nil {
// 		return false
// 	}

// 	for _, group := range uGroupList {
// 		if group == "tikiadmins" {
// 			return true
// 		}

// 	}
// 	return false
// }

// func getUserGroupMembershipGroupList(groupRepository *repositories.GroupRepository, userEmail string) ([]models.TicketGroup, error) {

// 	var userGroupsList []models.TicketGroup
// 	// search groups and group members

// 	existingGroups, err := groupRepository.GetAllGroups()
// 	//existingGroups, err := controller.GetAllGroupsOnDynamo()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, group := range existingGroups {
// 		gMembers := group.GroupMembers
// 		for _, member := range gMembers {
// 			if userEmail == member {
// 				userGroupsList = append(userGroupsList, group)
// 			}
// 		}

// 	}
// 	return userGroupsList, nil
// }

// func IsUserAllowedToPerformTicketOperation(groupRepository *repositories.GroupRepository, userEmail string, operationType string) bool {
// 	switch operationType {
// 	case "delete", "create", "setSecret", "getSecret":
// 		if isUserInTikiadmins(groupRepository, userEmail) {
// 			return true
// 		}

// 		userGroupMembershipGroupList, err := getUserGroupMembershipGroupList(groupRepository, userEmail)
// 		if err != nil {
// 			return false
// 		}

// 		for _, group := range userGroupMembershipGroupList {
// 			accPerms := group.AccessPerms
// 			if accPerms.Ticket[operationType] {
// 				return true
// 			}
// 		}
// 		return false
// 	case "show", "info", "list":
// 		return true
// 	default:
// 		return false
// 	}
// }

// func IsUserMemberOfTikiadmins(groupRepository *repositories.GroupRepository, userEmail string) bool {
// 	uGroupList, err := getUserGroupMembershipNames(groupRepository, userEmail)
// 	if err != nil {
// 		return false
// 	}

// 	for _, group := range uGroupList {
// 		if group == "tikiadmins" {
// 			return true
// 		}

// 	}
// 	return false
// }

//DoesGroupDomainScopeAllowedToUser
// func IsUserAllowedInGroupDomainScope(groupRepository *repositories.GroupRepository, userEmail string, ticketOrDomainPath string, domainScopeOperation string) bool {
// 	isUserAllowed := false

// 	ticketGroups, err := getUserGroupMembershipGroupList(groupRepository, userEmail)
// 	if err != nil {
// 		isUserAllowed = false
// 		return isUserAllowed
// 	}

// 	var domainScopeList []string

// 	for _, v := range ticketGroups {
// 		domainScopeList = append(domainScopeList, v.DomainScope.Root)
// 	}

// 	switch domainScopeOperation {
// 	case "createTicket", "deleteTicket", "assumeTicket":

// 		if IsUserMemberOfTikiadmins(groupRepository, userEmail) {
// 			return true
// 		}

// 		for _, v := range domainScopeList {
// 			if strings.HasPrefix(ticketOrDomainPath, v) {
// 				isUserAllowed = true
// 				return true
// 			}
// 		}

// 		return isUserAllowed
// 	case "createDomain", "deleteDomain":
// 		if IsUserMemberOfTikiadmins(groupRepository, userEmail) {
// 			return true
// 		}

// 		for _, v := range domainScopeList {
// 			if strings.HasPrefix(ticketOrDomainPath, v) {
// 				isUserAllowed = true
// 				return true
// 			}
// 		}
// 		return isUserAllowed

// 	default:
// 		return isUserAllowed
// 	}
// }
