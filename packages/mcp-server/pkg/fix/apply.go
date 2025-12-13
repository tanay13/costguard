package fix

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tanay13/costguard/packages/mcp-server/pkg/types"
)

func ApplyFix(action types.FixAction) error {
	switch action.Provider {
	case types.ProviderKubernetes:
		return applyK8sFix(action)
	default:
		return fmt.Errorf("unsupported provider: %s", action.Provider)
	}
}

func applyK8sFix(action types.FixAction) error {

	manifestFiles, err := findK8sManifests()
	if err != nil {
		return fmt.Errorf("failed to find manifests: %w", err)
	}

	if len(manifestFiles) == 0 {
		return fmt.Errorf("no Kubernetes manifest files found")
	}

	var targetFile string
	for _, file := range manifestFiles {
		if containsResource(file, action.Resource) {
			targetFile = file
			break
		}
	}

	if targetFile == "" {
		if len(action.FilesToEdit) > 0 {
			targetFile = action.FilesToEdit[0]
		} else {
			targetFile = manifestFiles[0]
		}
	}

	return updateK8sManifest(targetFile, action)
}

func findK8sManifests() ([]string, error) {
	var files []string

	searchPaths := []string{
		"k8s",
		"kubernetes",
		"deployments",
		"manifests",
		".",
	}

	for _, path := range searchPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if info.IsDir() {
				return nil
			}

			ext := filepath.Ext(p)
			if ext == ".yaml" || ext == ".yml" {

				content, err := os.ReadFile(p)
				if err != nil {
					return nil
				}

				if strings.Contains(string(content), "apiVersion:") &&
					(strings.Contains(string(content), "kind: Deployment") ||
						strings.Contains(string(content), "kind: StatefulSet") ||
						strings.Contains(string(content), "kind: Pod")) {
					files = append(files, p)
				}
			}

			return nil
		})

		if err != nil {
			continue
		}
	}

	return files, nil
}

func containsResource(filePath, resourceName string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}

	contentStr := string(content)
	return strings.Contains(contentStr, resourceName)
}

func updateK8sManifest(filePath string, action types.FixAction) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return updateK8sManifestRegex(filePath, content, action)
}

func updateK8sManifestRegex(filePath string, content []byte, action types.FixAction) error {
	contentStr := string(content)

	var pattern *regexp.Regexp
	var replacement string

	if action.Action.Field == "resources.requests.cpu" {

		pattern = regexp.MustCompile(`(?i)(\s+requests:\s*\n\s+)?cpu:\s*["']?(\d+(?:\.\d+)?)(m)?["']?`)
		replacement = fmt.Sprintf("${1}cpu: \"%.0fm\"", action.Action.Value)
	} else if action.Action.Field == "resources.requests.memory" {

		memGi := action.Action.Value
		pattern = regexp.MustCompile(`(?i)(\s+requests:\s*\n\s+)?memory:\s*["']?(\d+(?:\.\d+)?)(Mi|Gi|M|G)?["']?`)
		replacement = fmt.Sprintf("${1}memory: \"%.2fGi\"", memGi)
	} else {
		return fmt.Errorf("unsupported field: %s", action.Action.Field)
	}

	lines := strings.Split(contentStr, "\n")
	updated := false
	inResource := false
	resourceIndent := 0

	for i, line := range lines {

		if strings.Contains(line, "name:") && strings.Contains(line, action.Resource) {
			inResource = true
			resourceIndent = len(line) - len(strings.TrimLeft(line, " "))
			continue
		}

		if inResource {

			currentIndent := len(line) - len(strings.TrimLeft(line, " "))
			if currentIndent <= resourceIndent && strings.TrimSpace(line) != "" && !strings.HasPrefix(strings.TrimSpace(line), "#") {
				if !strings.Contains(line, "name:") {
					inResource = false
					continue
				}
			}

			if pattern.MatchString(line) {
				lines[i] = pattern.ReplaceAllString(line, replacement)
				updated = true
			}
		}
	}

	if !updated {

		return addResourcesSection(filePath, contentStr, action)
	}

	return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
}

func addResourcesSection(filePath, content string, action types.FixAction) error {

	lines := strings.Split(content, "\n")

	for i, line := range lines {
		if strings.Contains(line, "containers:") || strings.Contains(line, "- name:") {

			indent := strings.Repeat(" ", 6)

			insertIndex := i + 1
			resourcesLines := []string{
				indent + "resources:",
				indent + "  requests:",
			}

			if action.Action.Field == "resources.requests.cpu" {
				resourcesLines = append(resourcesLines, fmt.Sprintf("%s    cpu: \"%.0fm\"", indent, action.Action.Value))
			} else if action.Action.Field == "resources.requests.memory" {
				resourcesLines = append(resourcesLines, fmt.Sprintf("%s    memory: \"%.2fGi\"", indent, action.Action.Value))
			}

			newLines := make([]string, 0, len(lines)+len(resourcesLines))
			newLines = append(newLines, lines[:insertIndex]...)
			newLines = append(newLines, resourcesLines...)
			newLines = append(newLines, lines[insertIndex:]...)

			return os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644)
		}
	}

	return fmt.Errorf("could not find container section to add resources")
}

func updateK8sManifestYAML(filePath string, doc interface{}, action types.FixAction) error {
	content, _ := os.ReadFile(filePath)
	return updateK8sManifestRegex(filePath, content, action)
}
