package line

// FlexMessageBuilder defines interface for building flex messages
type FlexMessageBuilder interface {
	Build() map[string]interface{}
}

// FlexTemplateType represents types of flex templates
type FlexTemplateType string

const (
	TemplateMainMenu        FlexTemplateType = "main_menu"
	TemplateEquipmentChange FlexTemplateType = "equipment_change"
	TemplateContactStaff    FlexTemplateType = "contact_staff"
)

// Color constants for consistent branding
const (
	ColorPrimary   = "#0367D3"
	ColorSecondary = "#B8D4F0"
	ColorOrange    = "#FF9800"
	ColorGreen     = "#4CAF50"
	ColorBlue      = "#5B9BD5"
	ColorGray      = "#AAAAAA"
	ColorRed       = "#C62828"
	ColorRedLight  = "#FFEBEE"
	ColorWhite     = "#FFFFFF"
)
