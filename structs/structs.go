package structs

// ColorCatalog is root element of a Bricklink color xml list export
type ColorCatalog struct {
	Items []ColorItem `xml:"ITEM"`
}

// ColorItem is a single item in the Bricklink color xml list
type ColorItem struct {
	Color     int    `xml:"COLOR"`
	ColorName string `xml:"COLORNAME"`
}

// Inventory is root element of a Bricklink wanted list export
type Inventory struct {
	Items []Item `xml:"ITEM"`
}

// Item is a single item in the Bricklink wanted list
type Item struct {
	ItemType  string  `xml:"ITEMTYPE"`
	ItemID    string  `xml:"ITEMID"`
	Color     int     `xml:"COLOR"`
	Maxprice  float32 `xml:"MAXPRICE"`
	MinQty    int     `xml:"MINQTY"`
	Condition string  `xml:"CONDITION"`
	Notify    string  `xml:"NOTIFY"`
}

// LegoSet represents the single set
type LegoSet struct {
	ID          	int
	Name        	string
	Description 	string
	RequiredCount 	int
	FoundCount	int
}

// LegoPart represents the single part
type LegoPart struct {
	ID          int
	Partnumber  string
	Description string
	ColorID     int
	SetID       int
	RequiredQty int
	FoundQty    int
	LowPriority int
}

// LegoColor represents a color
type LegoColor struct {
	ID     int
	Number int
	Name   string
}
