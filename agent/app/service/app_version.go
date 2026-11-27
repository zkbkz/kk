package service

import (
	"regexp"
	"strings"

	"github.com/1Panel-dev/1Panel/agent/utils/common"
)

var appVersionPattern = regexp.MustCompile(`^[vV]?\d+(?:\.\d+)*(?:[-+][A-Za-z0-9][A-Za-z0-9.-]*)?$`)

func isNewerAppVersion(_ string, target, current string) bool {
	if !appVersionPattern.MatchString(target) || !appVersionPattern.MatchString(current) {
		return false
	}
	return common.CompareVersion(target, current)
}

func isCrossAppVersion(_ string, current, target string) bool {
	if !appVersionPattern.MatchString(current) || !appVersionPattern.MatchString(target) {
		return false
	}
	current = strings.TrimPrefix(strings.TrimPrefix(current, "v"), "V")
	target = strings.TrimPrefix(strings.TrimPrefix(target, "v"), "V")
	return common.IsCrossVersion(current, target)
}

func sortAppVersions(_ string, versions []string) []string {
	return common.GetSortedVersions(versions)
}
