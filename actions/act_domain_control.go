package actions

import "strings"

func UpdateParent(domainPath string) string {
	domainPath = TrimDomainPath(domainPath)
	if !strings.Contains(domainPath, "/") || domainPath == "/" {
		return "/"
	}

	pObjs := strings.Split(domainPath, "/")
	parent := strings.Join(pObjs[:len(pObjs)-1], "/")

	return parent
}

func TrimDomainPath(domainPath string) string {
	if domainPath == "/" {
		return domainPath
	}

	return strings.TrimRight(domainPath, "/")
}
