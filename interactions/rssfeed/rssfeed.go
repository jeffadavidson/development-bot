package rssfeed

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/jeffadavidson/development-bot/utilities/fileio"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
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
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

// CreateRSSFeed creates a new RSS feed with the given title and description
func CreateRSSFeed(title, description, link string) *RSS {
	return &RSS{
		Version: "2.0",
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
func (rss *RSS) AddItem(title, description, link, guid string, pubDate time.Time) {
	item := Item{
		Title:       title,
		Link:        link,
		Description: description,
		PubDate:     pubDate.Format(time.RFC1123Z),
		GUID:        guid,
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

// UpdateItem updates an existing item or adds it if not found
func (rss *RSS) UpdateItem(title, description, link, guid string, pubDate time.Time) {
	item := rss.FindItemByGUID(guid)
	if item != nil {
		// Update existing item
		item.Title = title
		item.Description = description
		item.Link = link
		item.PubDate = pubDate.Format(time.RFC1123Z)
	} else {
		// Add new item
		rss.AddItem(title, description, link, guid, pubDate)
	}
	rss.Channel.LastBuildDate = time.Now().Format(time.RFC1123Z)
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