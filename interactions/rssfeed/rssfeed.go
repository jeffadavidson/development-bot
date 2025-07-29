package rssfeed

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/jeffadavidson/development-bot/utilities/fileio"
)

type RSS struct {
	XMLName   xml.Name `xml:"rss"`
	Version   string   `xml:"version,attr"`
	ContentNS string   `xml:"xmlns:content,attr"`
	Channel   Channel  `xml:"channel"`
}

type Channel struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	Language      string `xml:"language"`
	LastBuildDate string `xml:"lastBuildDate"`
	Items         []Item `xml:"item"`
}

type Item struct {
	Title          string    `xml:"title"`
	Link           string    `xml:"link"`
	Description    CDataText `xml:"description"`
	ContentEncoded CDataText `xml:"content:encoded"`
	PubDate        string    `xml:"pubDate"`
	GUID           string    `xml:"guid"`
	Category       string    `xml:"category,omitempty"`
	Author         string    `xml:"author,omitempty"`
	Source         string    `xml:"source,omitempty"`
	Comments       string    `xml:"comments,omitempty"`
}

// CDataText wraps text in CDATA sections for proper RSS compatibility
type CDataText struct {
	Text string `xml:",cdata"`
}

// CreateRSSFeed creates a new RSS feed with the given title and description
func CreateRSSFeed(title, description, link string) *RSS {
	return &RSS{
		Version:   "2.0",
		ContentNS: "http://purl.org/rss/1.0/modules/content/",
		Channel: Channel{
			Title:         title,
			Link:          link,
			Description:   description,
			Language:      "en-us",
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			Items:         []Item{},
		},
	}
}

// AddItem adds a new item to the RSS feed
func (rss *RSS) AddItem(title, description, link, guid string, pubDate time.Time, category, author, source, comments, contentEncoded string) {
	item := Item{
		Title:          title,
		Link:           link,
		Description:    CDataText{Text: description},
		ContentEncoded: CDataText{Text: contentEncoded},
		PubDate:        pubDate.Format(time.RFC1123Z),
		GUID:           guid,
		Category:       category,
		Author:         author,
		Source:         source,
		Comments:       comments,
	}

	// Add to beginning of slice (newest first)
	rss.Channel.Items = append([]Item{item}, rss.Channel.Items...)
	rss.Channel.LastBuildDate = time.Now().Format(time.RFC1123Z)
}

// ToXML converts the RSS feed to XML bytes
func (rss *RSS) ToXML() ([]byte, error) {
	xmlData, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal RSS to XML: %v", err)
	}

	// Add XML declaration
	xmlString := xml.Header + string(xmlData)
	return []byte(xmlString), nil
}

// LoadRSSFromXML loads an RSS feed from XML bytes
func LoadRSSFromXML(xmlData []byte) (*RSS, error) {
	var rss RSS
	err := xml.Unmarshal(xmlData, &rss)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal RSS from XML: %v", err)
	}
	return &rss, nil
}

// FindItemByGUID finds an item in the RSS feed by its GUID
func (rss *RSS) FindItemByGUID(guid string) *Item {
	for i := range rss.Channel.Items {
		if rss.Channel.Items[i].GUID == guid {
			return &rss.Channel.Items[i]
		}
	}
	return nil
}

// UpdateItem updates an existing item or adds a new one if it doesn't exist
// Returns true if any actual changes were made to the feed
func (rss *RSS) UpdateItem(title, description, link, guid string, pubDate time.Time, category, author, source, comments, contentEncoded string) bool {
	item := rss.FindItemByGUID(guid)
	if item != nil {
		// Check if any content actually changed
		newItem := Item{
			Title:          title,
			Description:    CDataText{Text: description},
			ContentEncoded: CDataText{Text: contentEncoded},
			Link:           link,
			PubDate:        pubDate.Format(time.RFC1123Z),
			Category:       category,
			Author:         author,
			Source:         source,
			Comments:       comments,
		}

		// Compare current item with new item to detect changes
		hasChanges := item.Title != newItem.Title ||
			item.Description.Text != newItem.Description.Text ||
			item.ContentEncoded.Text != newItem.ContentEncoded.Text ||
			item.Link != newItem.Link ||
			item.PubDate != newItem.PubDate ||
			item.Category != newItem.Category ||
			item.Author != newItem.Author ||
			item.Source != newItem.Source ||
			item.Comments != newItem.Comments

		if hasChanges {
			// Update existing item
			item.Title = newItem.Title
			item.Description = newItem.Description
			item.ContentEncoded = newItem.ContentEncoded
			item.Link = newItem.Link
			item.PubDate = newItem.PubDate
			item.Category = newItem.Category
			item.Author = newItem.Author
			item.Source = newItem.Source
			item.Comments = newItem.Comments
			rss.Channel.LastBuildDate = time.Now().Format(time.RFC1123Z)
			return true
		}
		return false
	} else {
		// Add new item
		rss.AddItem(title, description, link, guid, pubDate, category, author, source, comments, contentEncoded)
		return true
	}
}

// TrimToMaxItems keeps only the most recent N items in the feed
func (rss *RSS) TrimToMaxItems(maxItems int) {
	if len(rss.Channel.Items) > maxItems {
		rss.Channel.Items = rss.Channel.Items[:maxItems]
	}
}

// GetOrCreateRSSFeed loads an existing RSS feed from file or creates a new one
func GetOrCreateRSSFeed(filepath, title, description, link string) (*RSS, error) {
	// Try to load existing feed
	xmlData, err := loadRSSFile(filepath)
	if err != nil {
		// Create new feed if file doesn't exist
		return CreateRSSFeed(title, description, link), nil
	}

	// Load existing feed
	rss, err := LoadRSSFromXML(xmlData)
	if err != nil {
		return nil, fmt.Errorf("failed to load existing RSS feed: %v", err)
	}

	// Ensure ContentNS is properly set (fix for existing feeds with empty namespace)
	if rss.ContentNS == "" {
		rss.ContentNS = "http://purl.org/rss/1.0/modules/content/"
	}

	// Fix empty content:encoded fields by copying from description
	for i := range rss.Channel.Items {
		if rss.Channel.Items[i].ContentEncoded.Text == "" && rss.Channel.Items[i].Description.Text != "" {
			rss.Channel.Items[i].ContentEncoded.Text = rss.Channel.Items[i].Description.Text
		}
	}

	return rss, nil
}

// SaveRSSFeed saves an RSS feed to a file
func SaveRSSFeed(rss *RSS, filepath string) error {
	xmlData, err := rss.ToXML()
	if err != nil {
		return fmt.Errorf("failed to convert RSS to XML: %v", err)
	}

	return saveRSSFile(filepath, xmlData)
}

// loadRSSFile reads RSS XML from file using the existing fileio utility
func loadRSSFile(filepath string) ([]byte, error) {
	return fileio.GetFileContents(filepath)
}

// saveRSSFile writes RSS XML to file using the existing fileio utility
func saveRSSFile(filepath string, data []byte) error {
	return fileio.WriteFileContents(filepath, data)
}
