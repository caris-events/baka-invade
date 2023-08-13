package main

import "github.com/samber/lo"

type FileType int

const (
	FileTypeUnknown FileType = iota
	FileTypeObject
	FileTypeObjectLogo
	FileTypeDict
)

type RemovalType string

const (
	RemovalTypeUnknown RemovalType = ""
	// 表示原作者已經來信要求移除圖片。
	RemovalTypeImage RemovalType = "image"
)

// Category
type Category struct {
	Code string `json:"code"`
	Text string `json:"text"`
}

// Cache
type Cache struct {
	Objects     map[string]string `json:"objects"`
	ObjectLogos map[string]string `json:"object_logos"`
	Dicts       map[string]string `json:"dicts"`
}

type Relations struct {
	Categories     []string          `json:"categories"`
	CategoryTexts  []string          `json:"category_texts"`
	Owners         []ObjectOwner     `json:"owners"`
	OwnerTexts     []string          `json:"owner_texts"`
	Types          []ObjectType      `json:"types"`
	TypeTexts      []string          `json:"type_texts"`
	Invasions      []ObjectInvasion  `json:"invasions"`
	InvasionTexts  []string          `json:"invasion_texts"`
	Directions     []ObjectDirection `json:"directions"`
	DirectionTexts []string          `json:"direction_texts"`
	Incomes        []ObjectIncome    `json:"incomes"`
	IncomeTexts    []string          `json:"income_texts"`
}

var (
	categories = []*Category{
		&Category{
			Code: "GP",
			Text: "集團",
		},
		&Category{
			Code: "FI",
			Text: "銀行/金融",
		},
		&Category{
			Code: "PP",
			Text: "房地產/能源",
		},
		&Category{
			Code: "PL",
			Text: "政治",
		},
		&Category{
			Code: "ME",
			Text: "傳播媒體",
		},
		&Category{
			Code: "BK",
			Text: "教育/書商/出版社",
		},
		&Category{
			Code: "SP",
			Text: "購物",
		},
		&Category{
			Code: "CY",
			Text: "日常用品",
		},
		&Category{
			Code: "EA",
			Text: "家具/電器",
		},
		&Category{
			Code: "MD",
			Text: "醫藥/美容",
		},
		&Category{
			Code: "FD",
			Text: "飲食",
		},
		&Category{
			Code: "ST",
			Text: "運動",
		},
		&Category{
			Code: "EN",
			Text: "娛樂",
		},
		&Category{
			Code: "GA",
			Text: "遊戲",
		},
		&Category{
			Code: "AP",
			Text: "軟體/網頁",
		},
		&Category{
			Code: "IT",
			Text: "資訊與科技",
		},
		&Category{
			Code: "MP",
			Text: "通訊商/手機",
		},
		&Category{
			Code: "PC",
			Text: "電腦/硬體",
		},
		&Category{
			Code: "TS",
			Text: "交通",
		},
		&Category{
			Code: "TL",
			Text: "旅行",
		},
	}
)

var relations = &Relations{
	Categories: lo.Map(categories, func(v *Category, _ int) string {
		return v.Code
	}),
	CategoryTexts: lo.Map(categories, func(v *Category, _ int) string {
		return v.Text
	}),
	Owners: []ObjectOwner{
		ObjectOwnerTaiwaneseGov,
		ObjectOwnerTaiwanese,
		ObjectOwnerTaiwaneseFake,
		ObjectOwnerHongKongese,
		ObjectOwnerHongKongeseFake,
		ObjectOwnerChineseGov,
		ObjectOwnerChinese,
		ObjectOwnerForeignGov,
		ObjectOwnerForeign,
		ObjectOwnerForeignFake,
	},
	OwnerTexts: []string{
		"台灣國有",
		"台灣",
		"中國在台",
		"香港",
		"中國",
		"中共國有",
		"中國",
		"外國國有",
		"外國",
		"中國在外",
	},
	Types: []ObjectType{
		ObjectTypePerson,
		ObjectTypeBrand,
		ObjectTypeCompany,
		ObjectTypeSoftware,
		ObjectTypeGame,
	},
	TypeTexts: []string{
		"人物",
		"品牌",
		"公司",
		"服務",
		"遊戲",
	},
	Invasions: []ObjectInvasion{
		ObjectInvasionSupport,
		ObjectInvasionIndirect,
		ObjectInvasionNeutral,
	},
	InvasionTexts: []string{
		"支持侵略",
		"間接協助",
		"未表態",
	},
	Directions: []ObjectDirection{
		ObjectDirectionChineseGov,
		ObjectDirectionChinese,
		ObjectDirectionNeutral,
	},
	DirectionTexts: []string{
		"跟隨中國政府",
		"討好中國市場",
		"試圖中立",
	},
	Incomes: []ObjectIncome{
		ObjectIncomeChineseGov,
		ObjectIncomeChineseShareholder,
		ObjectIncomeChinese,
	},
	IncomeTexts: []string{
		"中國政府",
		"中國股東",
		"中國市場",
	},
}
