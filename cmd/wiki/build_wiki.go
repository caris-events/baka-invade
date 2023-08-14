package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
	"github.com/samber/lo"
)

// PrepareObjects
func PrepareObjects(objects []*Object) {
	// Fill *GrandObject
	for _, v := range objects {
		if v.GrandCode == "" {
			continue
		}
		v.GrandObject = lo.FindOrElse(objects, nil, func(j *Object) bool {
			return j.Code == v.GrandCode
		})
	}

	//
	for _, v := range objects {
		v.RelatedObjects = RelatedObjects(objects, v)

		v.Categories = lo.Map(v.CategoryCodes, func(code string, _ int) *Category {
			return lo.FindOrElse(categories, nil, func(c *Category) bool {
				return code == c.Code
			})
		})
		v.RelatedBrands = lo.Filter(v.RelatedObjects, func(v *Object, _ int) bool {
			return v.Type == ObjectTypeBrand || v.Type == ObjectTypePerson || v.Type == ObjectTypeCompany
		})
		v.RelatedThings = lo.Filter(v.RelatedObjects, func(v *Object, _ int) bool {
			return v.Type == ObjectTypeSoftware || v.Type == ObjectTypeGame
		})
		v.Name = Name(v)
		v.SecondaryName = SecondaryName(v)
		v.OwnerStr = OwnerStr(v)
		v.Description = ParseContent(objects, v.Description)
		v.SubnameStr = strings.Join(v.Subnames, "、")
		v.TypeStr = TypeStr(v)

		buildObjectMessage(objects, v)

		this := v
		prev := v
		for {
			if this.GrandObject == nil {
				break
			}
			this = this.GrandObject

			v.Informations = append(v.Informations, &ObjectInformation{
				YearStr:     "—",
				Description: ParseContent(objects, fmt.Sprintf("[[%s]] 是 [[%s]] 的旗下%s。", prev.Code, this.Code, TypeStr(prev))),
			})
			v.Informations = append(v.Informations, this.Informations...)

			prev = this
		}

		//
		v.Informations = lo.UniqBy(v.Informations, func(v *ObjectInformation) string {
			return v.Description
		})

	}

	for _, v := range objects {
		var clearInvasion, clearDirection, clearIncome bool
		var lastInvasion ObjectInvasion
		var lastDirection ObjectDirection
		var lastIncome ObjectIncome
		for i := len(v.Informations) - 1; i >= 0; i-- {
			if clearIncome {
				v.Informations[i].Messages = lo.Filter(v.Informations[i].Messages, func(j *ObjectInformationMessage, _ int) bool {
					return !j.IsIncome
				})
			}
			if clearDirection {
				v.Informations[i].Messages = lo.Filter(v.Informations[i].Messages, func(j *ObjectInformationMessage, _ int) bool {
					return !j.IsDirection
				})
			}
			if clearInvasion {
				v.Informations[i].Messages = lo.Filter(v.Informations[i].Messages, func(j *ObjectInformationMessage, _ int) bool {
					return !j.IsInvasion
				})
			}
			if v.Informations[i].Income != "" && !clearIncome {
				lastIncome = v.Informations[i].Income
				clearIncome = true
			}
			if v.Informations[i].Direction != "" && !clearDirection {
				lastDirection = v.Informations[i].Direction
				clearDirection = true
			}
			if v.Informations[i].Invasion != "" && !clearInvasion {
				lastInvasion = v.Informations[i].Invasion
				clearInvasion = true
			}
		}

		v.Direction = lastDirection
		v.Income = lastIncome
		v.Invasion = lastInvasion
		v.DirectionStr = DirectionStr(v)
		v.InvasionStr = InvasionStr(v)
		v.IncomeStr = IncomeStr(v)
		v.Summary = SummaryStr(v)
	}
}

// SimpleObjects
func SimpleObjects(codes []string, objects []*Object) (simps []*SimpleObject) {
	for _, v := range objects {
		simps = append(simps, &SimpleObject{
			CodeID:      lo.IndexOf(codes, v.Code),
			GrandCodeID: lo.IndexOf(codes, v.GrandCode),
			CategoryCodeIDs: lo.Reduce(v.CategoryCodes, func(total []int, this string, _ int) []int {
				if i := lo.IndexOf(relations.Categories, this); i != -1 {
					total = append(total, i)
				}
				return total
			}, []int{}),
			SecondaryName: v.SecondaryName,
			Name:          v.Name,
			OwnerID:       lo.IndexOf(relations.Owners, v.Owner),
			TypeID:        lo.IndexOf(relations.Types, v.Type),
			InvasionID:    lo.IndexOf(relations.Invasions, v.Invasion),
			DirectionID:   lo.IndexOf(relations.Directions, v.Direction),
			IncomeID:      lo.IndexOf(relations.Incomes, v.Income),
			Object:        v,
		})
	}
	return
}

func RelatedObjects(objects []*Object, object *Object) (relates []*Object) {
	grands := []string{object.Code}

	for this := object; this.GrandObject != nil; this = this.GrandObject {
		grands = append(grands, this.GrandObject.Code)
	}
	for _, v := range objects {
		if lo.Contains(grands, v.Code) {
			relates = append(relates, v)
		}
		for this := v; this.GrandObject != nil; this = this.GrandObject {
			if lo.Contains(grands, this.GrandObject.Code) {
				relates = append(relates, v)
				break
			}
		}
	}

	return lo.Uniq(lo.Filter(relates, func(v *Object, _ int) bool {
		return v.Code != object.Code
	}))
}

// OwnerStr
func OwnerStr(v *Object) string {
	t := TypeStr(v)

	switch v.Owner {
	case ObjectOwnerTaiwaneseGov:
		return "台灣國有" + t
	case ObjectOwnerTaiwanese:
		return "台灣" + t
	case ObjectOwnerTaiwaneseFake:
		return "中國在台" + t
	case ObjectOwnerChineseGov:
		return "中共國有" + t
	case ObjectOwnerChinese:
		return "中國" + t
	case ObjectOwnerHongKongese:
		return "香港" + t
	case ObjectOwnerHongKongeseFake:
		return "中國" + t
	case ObjectOwnerForeignGov:
		return "外國國有" + t
	case ObjectOwnerForeign:
		return "外國" + t
	case ObjectOwnerForeignFake:
		return "中國在外" + t
	default:
		return t
	}
}

// SummaryStr
func SummaryStr(v *Object) string {
	content := fmt.Sprintf("《%s》", Name(v))

	switch v.Invasion {
	case ObjectInvasionSupport:
		content += "是支持中國侵略、統一的"
	case ObjectInvasionIndirect:
		content += "是間接協助中國侵略、統一的"
	default: // ObjectInvasionNeutral
		content += "是"
	}

	// 台灣公司
	content += fmt.Sprintf("%s，", OwnerStr(v))

	switch v.Income {
	case ObjectIncomeChineseGov:
		content += "由中共政府持有管控，"
	case ObjectIncomeChineseShareholder:
		content += "背後由中國股東金錢資助，"
	case ObjectIncomeChinese:
		content += "目前主要資金市場來源在中國，"
	default:
		content += "受眾主要來自本國或全球，"
	}
	switch v.Direction {
	case ObjectDirectionChineseGov:
		content += "主要跟隨中共國家的政策營運（如：言論審查、限娛令）。"
	case ObjectDirectionChinese:
		content += "為了達成金錢利益會盡力討好中國市場。"
	case ObjectDirectionNeutral:
		content += "營運方式目前保持政治中立。"
	default:
		if (v.Income == "" || v.Income == ObjectIncomeUntracked) && (v.Direction == "" || v.Direction == ObjectDirectionUntracked) && (v.Invasion == "" || v.Invasion == ObjectInvasionUntracked) {
			content += "營運方式自由且不受中國的管控。"
		} else {
			content += "目前依照自主方式營運。"
		}

	}
	return content
}

// YearStr
func YearStr(v string) string {
	if v == "" {
		return "—"
	}
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		t, err := time.Parse("2006", v)
		if err != nil {
			return "—"
		}
		return strconv.Itoa(t.Year()) + " 年"
	}
	return strconv.Itoa(t.Year()) + " 年"
}

// SecondaryName
func SecondaryName(v *Object) string {
	if v.FullnameZH != "" && v.FullnameEN != "" {
		return v.FullnameEN
	}
	return ""
}

// TypeStr
func TypeStr(v *Object) string {
	switch v.Type {
	case ObjectTypeBrand:
		return "品牌"
	case ObjectTypePerson:
		return "人物"
	case ObjectTypeCompany:
		return "公司"
	case ObjectTypeSoftware:
		return "服務"
	case ObjectTypeGame:
		return "遊戲"
	}
	log.Fatalf("object has unknown type: %s\n", v.Code)
	return ""
}

// IncomeStr
func IncomeStr(v *Object) string {
	switch v.Income {
	case ObjectIncomeChineseGov:
		return "中國政府"
	case ObjectIncomeChineseShareholder:
		return "中國股東"
	case ObjectIncomeChinese:
		return "中國市場"
	default:
		return "—"
	}
}

// DirectionStr
func DirectionStr(v *Object) string {
	switch v.Direction {
	case ObjectDirectionChineseGov:
		return "跟隨中國政府"
	case ObjectDirectionChinese:
		return "討好中國市場"
	case ObjectDirectionNeutral:
		return "試圖中立"
	default:
		return "—"
	}
}

// InvasionStr
func InvasionStr(v *Object) string {
	switch v.Invasion {
	case ObjectInvasionSupport:
		return "支持侵略"
	case ObjectInvasionIndirect:
		return "間接協助"
	case ObjectInvasionNeutral:
		return "未表態"
	default:
		return "—"
	}
}

// Name
func Name(v *Object) string {
	if v.FullnameZH != "" {
		return v.FullnameZH
	}
	if v.FullnameEN != "" {
		return v.FullnameEN
	}
	if len(v.Subnames) > 0 {
		return v.Subnames[0]
	}
	log.Fatalf("object has no name: %s\n", v.Code)
	return ""
}

// SelfStr
func SelfStr(v *Object) (output string) {
	switch v.Type {
	case ObjectTypeBrand, ObjectTypePerson, ObjectTypeSoftware, ObjectTypeGame:
		output += "這個" + TypeStr(v)
	case ObjectTypeCompany:
		output += "這間" + TypeStr(v)
	}
	return
}

// ParseContent
func ParseContent(objects []*Object, v string) string {
	v = string(blackfriday.Run([]byte(v)))

	// 將 [[hololive]] 轉成帶有連結的 Object["hololive"].FullnameZH
	v = ReplaceAllStringSubmatchFunc(regexp.MustCompile(` ?\[\[(.*?)\]\] ?`), v, func(groups []string) string {
		if len(groups) == 0 {
			return ""
		}
		object, ok := lo.Find(objects, func(v *Object) bool {
			return v.Code == groups[1]
		})
		if !ok {
			// log.Fatalf("object not found: %s\n", groups[1])
			return fmt.Sprintf(`《%s》`, groups[1])
		}
		return fmt.Sprintf(`《<a class="ts-text is-link" href="./../%s">%s</a>》`, object.Code, Name(object))
	})
	return v
}

// buildObjectMessage
func buildObjectMessage(objects []*Object, v *Object) {
	for _, info := range v.Informations {
		info.YearStr = YearStr(info.Date)
		info.Description = ParseContent(objects, info.Description)

		// Income
		switch info.Income {
		case ObjectIncomeChineseGov:
			info.Messages = append(info.Messages, &ObjectInformationMessage{
				IsIncome: true,
				Type:     "warning",
				Text:     SelfStr(v) + "已經被中國政府接管。",
			})
		case ObjectIncomeChineseShareholder:
			info.Messages = append(info.Messages, &ObjectInformationMessage{
				IsIncome: true,
				Type:     "warning",
				Text:     SelfStr(v) + "背後由中國股東金錢資助。",
			})
		case ObjectIncomeChinese:
			info.Messages = append(info.Messages, &ObjectInformationMessage{
				IsIncome: true,
				Type:     "warning",
				Text:     SelfStr(v) + "主要資金來源是中國與其市場。",
			})
		case ObjectIncomeUntracked:
			info.Messages = append(info.Messages, &ObjectInformationMessage{
				IsIncome: true,
				Type:     "info",
				Text:     SelfStr(v) + "的客戶市場與收入來自本國或世界各地，不會特別偏好中國市場。",
			})
		}

		// Direction
		switch info.Direction {
		case ObjectDirectionChineseGov:
			info.Messages = append(info.Messages, &ObjectInformationMessage{
				IsDirection: true,
				Type:        "warning",
				Text:        SelfStr(v) + "營運方針主要遵循中共政府的指示。",
			})
		case ObjectDirectionChinese:
			info.Messages = append(info.Messages, &ObjectInformationMessage{
				IsDirection: true,
				Type:        "warning",
				Text:        SelfStr(v) + "營運方針以討好中國市場為主。",
			})
		case ObjectDirectionNeutral:
			info.Messages = append(info.Messages, &ObjectInformationMessage{
				IsDirection: true,
				Type:        "info",
				Text:        SelfStr(v) + "營運方針盡可能地避免參與政治。",
			})
		case ObjectDirectionUntracked:
			info.Messages = append(info.Messages, &ObjectInformationMessage{
				IsDirection: true,
				Type:        "info",
				Text:        SelfStr(v) + "營運方針不會因為中國打壓而有所限制。",
			})
		}

		// Invasion
		switch info.Invasion {
		case ObjectInvasionSupport:
			info.Messages = append(info.Messages, &ObjectInformationMessage{
				IsInvasion: true,
				Type:       "warning",
				Text:       SelfStr(v) + "支持中國侵略、統一。",
			})
		case ObjectInvasionIndirect:
			info.Messages = append(info.Messages, &ObjectInformationMessage{
				IsInvasion: true,
				Type:       "warning",
				Text:       SelfStr(v) + "因為利益而間接協助了中國侵略、統一。",
			})
		case ObjectInvasionNeutral:
			info.Messages = append(info.Messages, &ObjectInformationMessage{
				IsInvasion: true,
				Type:       "info",
				Text:       SelfStr(v) + "試圖保持中立，盡可能不做出任何侵略性行為。",
			})
		case ObjectInvasionUntracked:
			info.Messages = append(info.Messages, &ObjectInformationMessage{
				IsInvasion: true,
				Type:       "info",
				Text:       SelfStr(v) + "自由且不受中國的管控。",
			})
		}
	}
}

func Minus(a, b int) int {
	return a - b
}

func IsDangerousIncome(v ObjectIncome) bool {
	return v == ObjectIncomeChineseGov || v == ObjectIncomeChineseShareholder || v == ObjectIncomeChinese
}

func IsDangerousDirection(v ObjectDirection) bool {
	return v == ObjectDirectionChineseGov || v == ObjectDirectionChinese
}

func IsDangerousInvasion(v ObjectInvasion) bool {
	return v == ObjectInvasionSupport || v == ObjectInvasionIndirect
}

func IsDangerousOwner(v ObjectOwner) bool {
	return v == ObjectOwnerTaiwaneseFake || v == ObjectOwnerHongKongeseFake || v == ObjectOwnerChineseGov || v == ObjectOwnerChinese || v == ObjectOwnerForeignFake
}
