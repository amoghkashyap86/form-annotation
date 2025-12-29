package annotation

import (
	"encoding/json"
	"os"
)

type FormAnnotation struct {
	FormMetadata FormMetadata `json:"form_metadata"`
	Pages        []Page       `json:"pages"`
	FieldGroups  []FieldGroup `json:"field_groups,omitempty"`
}

type FormMetadata struct {
	FormID    string   `json:"form_id"`
	FormName  string   `json:"form_name"`
	Year      int      `json:"year"`
	PageCount int      `json:"page_count"`
	PageSize  PageSize `json:"page_size"`
}

type PageSize struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"unit"`
}

type Page struct {
	PageNumber int     `json:"page_number"`
	Fields     []Field `json:"fields"`
}

type FieldType string

const (
	FieldTypeText      FieldType = "text"
	FieldTypeCurrency  FieldType = "currency"
	FieldTypeNumeric   FieldType = "numeric"
	FieldTypeCheckbox  FieldType = "checkbox"
	FieldTypeDate      FieldType = "date"
	FieldTypeSegmented FieldType = "segmented"
	FieldTypeSignature FieldType = "signature"
)

type DataType string

const (
	DataTypeString  DataType = "string"
	DataTypeDecimal DataType = "decimal"
	DataTypeInteger DataType = "integer"
	DataTypeBoolean DataType = "boolean"
	DataTypeDate    DataType = "date"
)

type Field struct {
	FieldID    string      `json:"field_id"`
	IRSLineRef string      `json:"irs_line_reference,omitempty"`
	FieldType  FieldType   `json:"field_type"`
	DataType   DataType    `json:"data_type"`
	Position   Position    `json:"position,omitempty"`
	Segments   []Segment   `json:"segments,omitempty"`
	Style      *TextStyle  `json:"style,omitempty"`
	CheckStyle *CheckStyle `json:"check_style,omitempty"`
	Formatting *Formatting `json:"formatting,omitempty"`
	Validation *Validation `json:"validation,omitempty"`
	GroupID    string      `json:"group_id,omitempty"`
	FieldValue string      `json:"field_value"`
}

type Position struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"unit"`
}

type Segment struct {
	Position Position `json:"position"`
	Length   int      `json:"length"`
}

type TextStyle struct {
	FontFamily    string  `json:"font_family,omitempty"`
	FontSize      int     `json:"font_size,omitempty"`
	FontWeight    string  `json:"font_weight,omitempty"`
	TextAlign     string  `json:"text_align,omitempty"`
	VerticalAlign string  `json:"vertical_align,omitempty"`
	Color         string  `json:"color,omitempty"`
	LetterSpacing float64 `json:"letter_spacing,omitempty"`
}

type CheckStyle struct {
	MarkType   string `json:"mark_type"`
	MarkSize   int    `json:"mark_size"`
	MarkWeight string `json:"mark_weight"`
}

type Formatting struct {
	DecimalPlaces  int    `json:"decimal_places,omitempty"`
	ShowCommas     bool   `json:"show_commas,omitempty"`
	NegativeFormat string `json:"negative_format,omitempty"`
	Prefix         string `json:"prefix,omitempty"`
	Suffix         string `json:"suffix,omitempty"`
	DateFormat     string `json:"date_format,omitempty"`
	TextTransform  string `json:"text_transform,omitempty"`
}

type Validation struct {
	Required  bool    `json:"required,omitempty"`
	Pattern   string  `json:"pattern,omitempty"`
	Min       float64 `json:"min,omitempty"`
	Max       float64 `json:"max,omitempty"`
	MinLength int     `json:"min_length,omitempty"`
	MaxLength int     `json:"max_length,omitempty"`
}

type FieldGroup struct {
	GroupID   string   `json:"group_id"`
	GroupType string   `json:"group_type"`
	FieldIDs  []string `json:"field_ids"`
}

func LoadFromFile(filepath string) (*FormAnnotation, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var annotation FormAnnotation
	if err := json.Unmarshal(data, &annotation); err != nil {
		return nil, err
	}
	return &annotation, nil
}

// SaveToFile writes the FormAnnotation to a JSON file.
func (fa *FormAnnotation) SaveToFile(filepath string) error {
	data, err := json.MarshalIndent(fa, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, data, 0644)
}

// GetFieldByID finds a field by its ID across all pages.
func (fa *FormAnnotation) GetFieldByID(fieldID string) *Field {
	for i := range fa.Pages {
		for j := range fa.Pages[i].Fields {
			if fa.Pages[i].Fields[j].FieldID == fieldID {
				return &fa.Pages[i].Fields[j]
			}
		}
	}
	return nil
}

// GetFieldsByFieldValue finds all fields that match a field value path.
func (fa *FormAnnotation) GetFieldsByFieldValue(fieldValue string) []Field {
	var fields []Field
	for _, page := range fa.Pages {
		for _, field := range page.Fields {
			if field.FieldValue == fieldValue {
				fields = append(fields, field)
			}
		}
	}
	return fields
}

// GetFieldsOnPage returns all fields on a specific page.
func (fa *FormAnnotation) GetFieldsOnPage(pageNum int) []Field {
	for _, page := range fa.Pages {
		if page.PageNumber == pageNum {
			return page.Fields
		}
	}
	return nil
}

// GetFieldsByGroupID returns all fields belonging to a specific group.
func (fa *FormAnnotation) GetFieldsByGroupID(groupID string) []Field {
	var fields []Field
	for _, page := range fa.Pages {
		for _, field := range page.Fields {
			if field.GroupID == groupID {
				fields = append(fields, field)
			}
		}
	}
	return fields
}

// GetAllFields returns all fields across all pages.
func (fa *FormAnnotation) GetAllFields() []Field {
	var fields []Field
	for _, page := range fa.Pages {
		fields = append(fields, page.Fields...)
	}
	return fields
}

// ToJSON converts the FormAnnotation to a JSON string.
func (fa *FormAnnotation) ToJSON() (string, error) {
	data, err := json.MarshalIndent(fa, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON parses a JSON string into a FormAnnotation.
func FromJSON(jsonStr string) (*FormAnnotation, error) {
	var annotation FormAnnotation
	if err := json.Unmarshal([]byte(jsonStr), &annotation); err != nil {
		return nil, err
	}
	return &annotation, nil
}
