package entity

type Settings struct {
	Id                             uint `gorm:"column:id;primaryKey"`
	ClipboardImageDetectionEnabled bool `gorm:"column:clipboard_image_detection_enabled;not null;default:true"`
}

func (Settings) TableName() string {
	return "settings"
}
