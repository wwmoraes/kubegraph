package kubegraph

import "github.com/wwmoraes/kubegraph/icons"

// RestoreIcons expands the icon files into the target directory
func RestoreIcons(dir string) error {
	return icons.RestoreAssets(dir, "icons")
}
