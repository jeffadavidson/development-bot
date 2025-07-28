package rssfeed

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRSSFeed(t *testing.T) {
	// Test creating a new RSS feed
	title := "Test Feed"
	description := "Test Description"
	link := "https://example.com"

	rss := CreateRSSFeed(title, description, link)

	assert.Equal(t, "2.0", rss.Version)
	assert.Equal(t, title, rss.Channel.Title)
	assert.Equal(t, description, rss.Channel.Description)
	assert.Equal(t, link, rss.Channel.Link)
	assert.Equal(t, "en-us", rss.Channel.Language)
	assert.Empty(t, rss.Channel.Items)
}

func TestAddItem(t *testing.T) {
	rss := CreateRSSFeed("Test", "Test", "https://example.com")
	
	pubDate := time.Date(2025, 7, 28, 12, 0, 0, 0, time.UTC)
	
	rss.AddItem(
		"Test Item",
		"Test Description", 
		"https://example.com/item1",
		"guid-123",
		pubDate,
		"Development Permit",
		"Test Author",
		"Test Source",
		"https://example.com/comments",
	)

	require.Len(t, rss.Channel.Items, 1)
	
	item := rss.Channel.Items[0]
	assert.Equal(t, "Test Item", item.Title)
	assert.Equal(t, "Test Description", item.Description)
	assert.Equal(t, "https://example.com/item1", item.Link)
	assert.Equal(t, "guid-123", item.GUID)
	assert.Equal(t, "Development Permit", item.Category)
	assert.Equal(t, "Test Author", item.Author)
	assert.Equal(t, "Test Source", item.Source)
	assert.Equal(t, "https://example.com/comments", item.Comments)
}

func TestFindItemByGUID(t *testing.T) {
	rss := CreateRSSFeed("Test", "Test", "https://example.com")
	pubDate := time.Now()
	
	// Add test items
	rss.AddItem("Item 1", "Desc 1", "https://example.com/1", "guid-1", pubDate, "", "", "", "")
	rss.AddItem("Item 2", "Desc 2", "https://example.com/2", "guid-2", pubDate, "", "", "", "")
	
	// Test finding existing item
	foundItem := rss.FindItemByGUID("guid-1")
	require.NotNil(t, foundItem)
	assert.Equal(t, "Item 1", foundItem.Title)
	
	// Test not finding non-existent item
	notFoundItem := rss.FindItemByGUID("guid-nonexistent")
	assert.Nil(t, notFoundItem)
}

func TestUpdateItem_ExistingItem(t *testing.T) {
	rss := CreateRSSFeed("Test", "Test", "https://example.com")
	pubDate := time.Now()
	
	// Add initial item
	rss.AddItem("Old Title", "Old Desc", "https://example.com/old", "guid-123", pubDate, "Old Category", "", "", "")
	
	// Update the item
	newPubDate := pubDate.Add(time.Hour)
	rss.UpdateItem("New Title", "New Desc", "https://example.com/new", "guid-123", newPubDate, "New Category", "New Author", "", "")
	
	// Should still have only one item
	require.Len(t, rss.Channel.Items, 1)
	
	item := rss.Channel.Items[0]
	assert.Equal(t, "New Title", item.Title)
	assert.Equal(t, "New Desc", item.Description)
	assert.Equal(t, "https://example.com/new", item.Link)
	assert.Equal(t, "guid-123", item.GUID)
	assert.Equal(t, "New Category", item.Category)
	assert.Equal(t, "New Author", item.Author)
}

func TestUpdateItem_NewItem(t *testing.T) {
	rss := CreateRSSFeed("Test", "Test", "https://example.com")
	pubDate := time.Now()
	
	// Update non-existent item (should add it)
	rss.UpdateItem("New Item", "New Desc", "https://example.com/new", "guid-456", pubDate, "Category", "", "", "")
	
	require.Len(t, rss.Channel.Items, 1)
	
	item := rss.Channel.Items[0]
	assert.Equal(t, "New Item", item.Title)
	assert.Equal(t, "guid-456", item.GUID)
}

func TestTrimToMaxItems(t *testing.T) {
	rss := CreateRSSFeed("Test", "Test", "https://example.com")
	pubDate := time.Now()
	
	// Add 5 items
	for i := 0; i < 5; i++ {
		rss.AddItem(
			"Item "+string(rune('A'+i)),
			"Desc",
			"https://example.com/"+string(rune('A'+i)),
			"guid-"+string(rune('A'+i)),
			pubDate,
			"", "", "", "",
		)
	}
	
	require.Len(t, rss.Channel.Items, 5)
	
	// Trim to 3 items
	rss.TrimToMaxItems(3)
	
	assert.Len(t, rss.Channel.Items, 3)
	// Should keep most recent items (added to beginning)
	assert.Equal(t, "Item E", rss.Channel.Items[0].Title)
	assert.Equal(t, "Item D", rss.Channel.Items[1].Title)
	assert.Equal(t, "Item C", rss.Channel.Items[2].Title)
}

func TestToXML(t *testing.T) {
	rss := CreateRSSFeed("Test Feed", "Test Description", "https://example.com")
	pubDate := time.Date(2025, 7, 28, 12, 0, 0, 0, time.UTC)
	
	rss.AddItem(
		"Test Item",
		"Test Description",
		"https://example.com/item",
		"guid-123",
		pubDate,
		"Development Permit",
		"Test Author",
		"Test Source",
		"https://example.com/comments",
	)
	
	xmlData, err := rss.ToXML()
	require.NoError(t, err)
	
	xmlString := string(xmlData)
	
	// Check XML structure
	assert.Contains(t, xmlString, `<?xml version="1.0" encoding="UTF-8"?>`)
	assert.Contains(t, xmlString, `<rss version="2.0">`)
	assert.Contains(t, xmlString, `<title>Test Feed</title>`)
	assert.Contains(t, xmlString, `<description>Test Description</description>`)
	assert.Contains(t, xmlString, `<title>Test Item</title>`)
	assert.Contains(t, xmlString, `<guid>guid-123</guid>`)
	assert.Contains(t, xmlString, `<category>Development Permit</category>`)
	assert.Contains(t, xmlString, `<author>Test Author</author>`)
}

func TestLoadRSSFromXML(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Test Feed</title>
    <link>https://example.com</link>
    <description>Test Description</description>
    <language>en-us</language>
    <lastBuildDate>Mon, 28 Jul 2025 12:00:00 +0000</lastBuildDate>
    <item>
      <title>Test Item</title>
      <link>https://example.com/item</link>
      <description>Item Description</description>
      <pubDate>Mon, 28 Jul 2025 12:00:00 +0000</pubDate>
      <guid>guid-123</guid>
      <category>Development Permit</category>
    </item>
  </channel>
</rss>`

	rss, err := LoadRSSFromXML([]byte(xmlData))
	require.NoError(t, err)
	
	assert.Equal(t, "2.0", rss.Version)
	assert.Equal(t, "Test Feed", rss.Channel.Title)
	assert.Equal(t, "https://example.com", rss.Channel.Link)
	assert.Equal(t, "Test Description", rss.Channel.Description)
	
	require.Len(t, rss.Channel.Items, 1)
	item := rss.Channel.Items[0]
	assert.Equal(t, "Test Item", item.Title)
	assert.Equal(t, "guid-123", item.GUID)
	assert.Equal(t, "Development Permit", item.Category)
}

func TestLoadRSSFromXML_InvalidXML(t *testing.T) {
	invalidXML := `<invalid>xml</invalid>`
	
	_, err := LoadRSSFromXML([]byte(invalidXML))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal RSS from XML")
} 