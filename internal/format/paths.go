package format

// BuildCollectionPath generates a file path for a collection item.
// Returns: {collectionID}_{itemID}{extension}
func BuildCollectionPath(collectionID, itemID, extension string) string {
	return collectionID + "_" + itemID + extension
}

// BuildStandalonePath generates a file path for a standalone resource.
// Returns: {resourceID}{extension}
func BuildStandalonePath(resourceID, extension string) string {
	return resourceID + extension
}

// BuildClaudeCollectionPath generates a directory path for a Claude collection item.
// Returns: {collectionID}_{itemID}/SKILL.md
func BuildClaudeCollectionPath(collectionID, itemID string) string {
	return collectionID + "_" + itemID + "/SKILL.md"
}

// BuildClaudeStandalonePath generates a directory path for a Claude standalone resource.
// Returns: {resourceID}/SKILL.md
func BuildClaudeStandalonePath(resourceID string) string {
	return resourceID + "/SKILL.md"
}
