package notion

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestClient_ListDatabases_Integration(t *testing.T) {
	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		t.Skip("set NOTION_TOKEN to run this test")
	}

	wantTitle := "Task List 5132beee"

	s := New(token)

	// Get the list of the databases, list them one-by-one to exercise the pagination code path
	var got []Database

	page := Pagination{PageSize: 1}
	for {
		result, err := s.ListDatabases(context.Background(), page)
		if err != nil {
			t.Errorf("ListDatabases(%v) error = %v", page, err)
			return
		}
		for _, db := range result.Results {
			got = append(got, db)
		}
		if len(result.Results) > page.PageSize {
			t.Errorf("ListDatabases(%v) got too many items [%s]", page, renderTitles(result.Results))
			return
		}
		if !result.HasMore {
			break
		}
		page.StartCursor = result.NextCursor
	}

	// Check if there are at least two databases
	if len(got) < 2 {
		t.Errorf("Expected at least 2 databases, got [%s]", renderTitles(got))
		return
	}

	// Check if it contains the one we know it has to contain
	var taskList *Database
FindDB:
	for _, db := range got {
		for _, titlet := range db.Title {
			if strings.Contains(titlet.PlainText, wantTitle) {
				taskList = &db
				break FindDB
			}
		}
	}

	// If not then lets print what we have to make it easier to debug
	if taskList == nil {
		t.Errorf("Test DB [%s] not found. Got: [%s]", wantTitle, renderTitles(got))
	}
}

func renderTitles(got []Database) string {
	allTitles := []string{}
	for _, db := range got {
		title := []string{}
		for _, titlet := range db.Title {
			title = append(title, titlet.PlainText)
		}
		allTitles = append(allTitles, strings.Join(title, ""))
	}
	return strings.Join(allTitles, ", ")
}
