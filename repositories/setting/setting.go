package repo

type (
	// Setting store data
	Setting struct {
		Crud *SettingCRUD
	}
)

// NewSetting collection
func NewSetting() *Setting {
	return &Setting{Crud: NewCRUD()}
}
