package main

// ObjectType
type ObjectType string

const (
	ObjectTypeUnknown  ObjectType = ""
	ObjectTypePerson   ObjectType = "person"
	ObjectTypeBrand    ObjectType = "brand"
	ObjectTypeCompany  ObjectType = "company"
	ObjectTypeSoftware ObjectType = "software"
	ObjectTypeGame     ObjectType = "game"
	ObjectTypeGroup    ObjectType = "group"
)

// ObjectOwner
type ObjectOwner string

const (
	ObjectOwnerUnknown         ObjectOwner = ""
	ObjectOwnerTaiwaneseGov    ObjectOwner = "taiwanese_gov"
	ObjectOwnerTaiwanese       ObjectOwner = "taiwanese"
	ObjectOwnerTaiwaneseFake   ObjectOwner = "taiwanese_fake"
	ObjectOwnerHongKongese     ObjectOwner = "hongkongese"
	ObjectOwnerHongKongeseFake ObjectOwner = "hongkongese_fake"
	ObjectOwnerChineseGov      ObjectOwner = "chinese_gov"
	ObjectOwnerChinese         ObjectOwner = "chinese"
	ObjectOwnerForeignGov      ObjectOwner = "foreign_gov"
	ObjectOwnerForeign         ObjectOwner = "foreign"
	ObjectOwnerForeignFake     ObjectOwner = "foreign_fake"
)

// ObjectIncome
type ObjectIncome string

const (
	ObjectIncomeUnknown            ObjectIncome = ""
	ObjectIncomeChineseGov         ObjectIncome = "chinese_gov"
	ObjectIncomeChineseShareholder ObjectIncome = "chinese_shareholder"
	ObjectIncomeChinese            ObjectIncome = "chinese"
	ObjectIncomeUntracked          ObjectIncome = "untracked"
)

// ObjectDirection
type ObjectDirection string

const (
	ObjectDirectionUnknown    ObjectDirection = ""
	ObjectDirectionChineseGov ObjectDirection = "chinese_gov"
	ObjectDirectionChinese    ObjectDirection = "chinese"
	ObjectDirectionNeutral    ObjectDirection = "neutral"
	ObjectDirectionUntracked  ObjectDirection = "untracked"
)

// ObjectInvasion
type ObjectInvasion string

const (
	ObjectInvasionUnknown   ObjectInvasion = ""
	ObjectInvasionSupport   ObjectInvasion = "support"
	ObjectInvasionIndirect  ObjectInvasion = "indirect"
	ObjectInvasionNeutral   ObjectInvasion = "neutral"
	ObjectInvasionUntracked ObjectInvasion = "untracked"
)

// Object
type Object struct {
	Code          string               `yaml:"code"`
	GrandCode     string               `yaml:"grand_code"`
	FullnameZH    string               `yaml:"fullname_zh"`
	FullnameEN    string               `yaml:"fullname_en"`
	Subnames      []string             `yaml:"subnames"`
	Website       string               `yaml:"website"`
	Description   string               `yaml:"description"`
	CategoryCodes []string             `yaml:"categories"`
	Type          ObjectType           `yaml:"type"`
	Owner         ObjectOwner          `yaml:"owner"`
	Income        ObjectIncome         `yaml:"income"`
	Direction     ObjectDirection      `yaml:"direction"`
	Invasion      ObjectInvasion       `yaml:"invasion"`
	Informations  []*ObjectInformation `yaml:"informations"`

	Removal RemovalType `yaml:"removal"`

	Name          string `yaml:"-"`
	SecondaryName string `yaml:"-"`
	OwnerStr      string `yaml:"-"`
	DirectionStr  string `yaml:"-"`
	InvasionStr   string `yaml:"-"`
	IncomeStr     string `yaml:"-"`
	TypeStr       string `yaml:"-"`
	Summary       string `yaml:"-"`
	SubnameStr    string `yaml:"-"`

	Categories     []*Category `yaml:"-"`
	GrandObject    *Object     `yaml:"-"`
	ChildObjects   []*Object   `yaml:"-"`
	RelatedObjects []*Object   `yaml:"-"`
	RelatedBrands  []*Object   `yaml:"-"`
	RelatedThings  []*Object   `yaml:"-"`
}

// ObjectInformation
type ObjectInformation struct {
	Income      ObjectIncome    `yaml:"income"`
	Direction   ObjectDirection `yaml:"direction"`
	Invasion    ObjectInvasion  `yaml:"invasion"`
	Date        string          `yaml:"date"`
	Description string          `yaml:"description"`

	Messages []*ObjectInformationMessage
	YearStr  string
}

// ObjectInformationMessage
type ObjectInformationMessage struct {
	IsIncome    bool
	IsDirection bool
	IsInvasion  bool
	Type        string
	Text        string
}

// SimpleObject
type SimpleObject struct {
	CodeID          int    `json:"c"`
	GrandCodeID     int    `json:"g"`
	CategoryCodeIDs []int  `json:"r"`
	SecondaryName   string `json:"s"`
	Name            string `json:"n"`
	OwnerID         int    `json:"o"`
	TypeID          int    `json:"t"`
	InvasionID      int    `json:"i"`
	DirectionID     int    `json:"d"`
	IncomeID        int    `json:"m"`

	Code   string  `json:"code,omitempty"`
	Object *Object `json:"-"`
}
