// types for the bfpd database model
package bfpd

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Food struct {
	Id                int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         time.Time
	Ndbno             string `json:"ndbno" binding:"required" gorm:"unique;not null"`
	Description       string `json:"name" binding:"required"`
	ShortDescription  string
	CommercialName    string
	refuseDescription string
	FoodGroup         FoodGroup `gorm:"ForeignKey:FoodGroupID" json:"fg"`
	FoodGroupID       int64
	Ingredients       Ingredients `gorm:"ForeignKey:IngredientsID"`
	IngredientsID     int64
	Manufacturer      Manufacturer `json:"manu" gorm:"ForeignKey:ManufacturerID"`
	ManufacturerID    int64
	Datasource        string `json:"source"`
	Refuse            uint32
	ScientificName    string
	NFactor           float32
	ProFactor         float32
	FatFactor         float32
	ChoFactor         float32
	FootNote          []FootNote     `json:"footnotes" gorm:"ForeignKey:FoodID"`
	NutrientData      []NutrientData `json:"nutrients" gorm:"ForeignKey:FoodID"`
	Measures          []Weights      `json:"measures" gorm:"ForeignKey:FoodID"`
}
type FootNote struct {
	Id         int64
	version    uint32
	FnId       string   `json:"id" gorm:"not null"`
	FnType     string   `json:"type" gorm:"not null"`
	FnText     string   `json:"text" gorm:"not null"`
	Nutrient   Nutrient `gorm:"ForeignKey:NutrientID"`
	NutrientID int64    `json:"nutno"`
	FoodID     int64
}
type Ingredients struct {
	Id           int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
	Description  string `json:"desc" binding:"required"`
	Available    string `json:"avail"`        //	The date the data for the food item represented by the specific GTIN was made available on the market.
	Discontinued string `json:"discontinued"` //	The data indicated by the manufacturer that the product represented by a specific GTIN has been discontinued
	Updated      string `json:"updated"`
}
type Weights struct {
	Id           int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
	Version      uint8
	Seq          uint8   `json:"seq"`
	Amount       float32 `json:"amt"`
	Description  string  `json:"unit"`
	Gramweight   float32 `json:"weight"`
	Datapoints   uint32  `json:"dp"`
	Stddeviation float32 `json:'sdv'`
	FoodID       int64
}
type Nutrient struct {
	Id           int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
	Version      uint8
	Nutrientno   uint   `json:"nutno" binding:"required" gorm:"unique;not null"`
	Tagname      string `json:"tag"`
	Description  string `json:"desc" gorm:"not null"`
	Decimalpoint uint8  `json:"dp"`
	Srnutorder   uint32 `json:"sort"`
	Unit         Unit   `gorm:"ForeignKey:UnitID"`
	UnitID       uint
	Type         string

	NutrientData []NutrientData `gorm:"ForeignKey:NutrientID"`
}
type Manufacturer struct {
	Id        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Version   uint8
	Name      string `json:"name" binding:"required"`
	Foods     []Food
}

type Unit struct {
	Id        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Version   uint8  `json:"version"`
	Unit      string `gorm:"unique;not null" json:"unit"`
	Nutrients []Nutrient
}
type SourceCode struct {
	Id           int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
	Code         string         `binding:"required" json:"code"`
	Description  string         `json:"desc"`
	NutrientData []NutrientData //`gorm:"ForeignKey:SourceID"`
}
type Derivation struct {
	Id           int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
	Code         string `binding:"required" json:"code"`
	Description  string `json:"desc"`
	NutrientData []NutrientData
}
type FoodGroup struct {
	Id          int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	Cd          string `json:"cd" gorm:"unique;not null"`
	Description string `json:"desc"`
	Food        []Food
}
type NutrientData struct {
	Id             int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
	Value          float32 `json:"value"`
	Datapoints     uint32  `json:"dp"`
	StandardError  float32 `json:"se"`
	AddNutMark     string
	NumberStudies  uint8
	Minimum        float32
	Maximum        float32
	DegreesFreedom float32
	LowerEB        float32
	UpperEB        float32
	Comment        string
	ConfidenceCode string
	Sourcecode     SourceCode `json:"source" gorm:"ForeignKey:SourceID"`
	SourceID       int64
	Derivation     Derivation `json:"deriv" gorm:"ForeignKey:DerivationID"`
	DerivationID   int64
	Nutrient       Nutrient `gorm:"ForeignKey:NutrientID"`
	NutrientID     int64    `json:"nutno"`
	Food           Food     `gorm:"ForeignKey:FoodID"`
	FoodID         int64
}

// Database configuration
type Config struct {
	Db   string
	User string
	Pw   string
}
type DB struct {
	*gorm.DB
}
