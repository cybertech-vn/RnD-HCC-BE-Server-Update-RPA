package model

type SoftwareVersion struct {
	ID        string `json:"id" gorm:"primaryKey;column:id"`
	AppID     string `json:"app_id" gorm:"column:app_id;index"`
	Version   string `json:"version" gorm:"column:version;index"`
	FileName  string `json:"filename" gorm:"column:filename"`
	Checksum  string `json:"checksum" gorm:"column:checksum"`
	Size      int64  `json:"size" gorm:"column:size"`
	CreatedAt int64  `json:"created_at" gorm:"column:created_at;index"`
}

func (SoftwareVersion) TableName() string {
	return "software_versions"
}
