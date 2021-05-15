package notion

import (
	"os"
	"strings"
	"testing"
)

func TestClient_ListDatabases_Integration(t *testing.T) {
	fooAddr := os.Getenv("NOTION_TOKEN")
	if fooAddr == "" {
		t.Skip("set NOTION_TOKEN to run this test")
	}

	wantTitle := "Task List 5132beee"

	// Get the list of the databases
	c := &Client{Token: os.Getenv("NOTION_TOKEN")}
	got, err := c.ListDatabases()
	if err != nil {
		t.Errorf("ListDatabases() error = %v", err)
		return
	}

	// Check if it contains the one we know it has to contain
	var taskList *Database
FindDB:
	for _, db := range got.Results {
		for _, titlet := range db.Title {
			if strings.Contains(titlet.PlainText, wantTitle) {
				taskList = &db
				break FindDB
			}
		}
	}

	// If not then lets print what we have to make it easier to debug
	if taskList == nil {
		allTitles := []string{}
		for _, db := range got.Results {
			title := []string{}
			for _, titlet := range db.Title {
				title = append(title, titlet.PlainText)
			}
			allTitles = append(allTitles, strings.Join(title, ""))
		}
		t.Errorf("Test DB [%s] not found. Got: [%s]", wantTitle, strings.Join(allTitles, ", "))
	}
}
