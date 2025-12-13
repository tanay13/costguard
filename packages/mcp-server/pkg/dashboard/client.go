package dashboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/tanay13/costguard/packages/mcp-server/pkg/types"
)

type DashboardClient struct {
	BaseURL string
}

type DashboardUpdate struct {
	RepoOwner    string                  `json:"repo_owner"`
	RepoName     string                  `json:"repo_name"`
	RepoFullName string                  `json:"repo_full_name"`
	ScanData     types.ScanResponse      `json:"scan_data"`
	DecisionData types.AIDecisionSummary `json:"decision_data"`
	PRURL        string                  `json:"pr_url,omitempty"`
	PRNumber     int                     `json:"pr_number,omitempty"`
	Timestamp    string                  `json:"timestamp"`
}

func NewClient(baseURL string) *DashboardClient {
	return &DashboardClient{BaseURL: baseURL}
}

func GetRepoInfo() (owner, name, fullName string, err error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get git remote: %w", err)
	}

	remoteURL := strings.TrimSpace(string(output))

	if strings.HasPrefix(remoteURL, "git@github.com:") {
		remoteURL = strings.TrimPrefix(remoteURL, "git@github.com:")
		remoteURL = strings.TrimSuffix(remoteURL, ".git")
		parts := strings.Split(remoteURL, "/")
		if len(parts) >= 2 {
			owner = parts[0]
			name = parts[1]
			fullName = owner + "/" + name
			return owner, name, fullName, nil
		}
	}

	// Handle https://github.com/owner/repo.git format
	re := regexp.MustCompile(`github\.com[/:]([^/]+)/([^/]+?)(?:\.git)?$`)
	matches := re.FindStringSubmatch(remoteURL)
	if len(matches) >= 3 {
		owner = matches[1]
		name = matches[2]
		fullName = owner + "/" + name
		return owner, name, fullName, nil
	}

	return "", "", "", fmt.Errorf("could not parse repo info from: %s", remoteURL)
}

func (c *DashboardClient) SendUpdate(scanData types.ScanResponse, decisionData types.AIDecisionSummary, prURL string, prNumber int) error {
	if c.BaseURL == "" {

		return nil
	}

	owner, name, fullName, err := GetRepoInfo()
	if err != nil {
		// If we can't get repo info, still try to send (with empty repo fields)
		owner = "unknown"
		name = "unknown"
		fullName = "unknown/unknown"
	}

	update := DashboardUpdate{
		RepoOwner:    owner,
		RepoName:     name,
		RepoFullName: fullName,
		ScanData:     scanData,
		DecisionData: decisionData,
		PRURL:        prURL,
		PRNumber:     prNumber,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal update: %w", err)
	}

	url := strings.TrimSuffix(c.BaseURL, "/") + "/api/submit"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if apiKey := os.Getenv("COSTGUARD_DASHBOARD_API_KEY"); apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("dashboard returned status %d", resp.StatusCode)
	}

	return nil
}
