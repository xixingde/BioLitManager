package entity

import (
	"time"

	"gorm.io/gorm"
)

// Paper 论文实体
type Paper struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	Title                 string         `gorm:"type:varchar(255);not null" json:"title"`
	Abstract              string         `gorm:"type:text" json:"abstract"`
	JournalID             uint           `gorm:"index" json:"journal_id"`
	DOI                   string         `gorm:"type:varchar(100)" json:"doi"`
	PubMedID              string         `gorm:"type:varchar(50)" json:"pubmed_id"`
	ISSN                  string         `gorm:"type:varchar(20)" json:"issn"`
	ImpactFactor          float64        `json:"impact_factor"`
	Partition             string         `gorm:"type:varchar(50)" json:"partition"`
	IsSCI                 bool           `gorm:"default:false" json:"is_sci"`
	IsEI                  bool           `gorm:"default:false" json:"is_ei"`
	IsCI                  bool           `gorm:"default:false" json:"is_ci"`
	IsDI                  bool           `gorm:"default:false" json:"is_di"`
	IsCore                bool           `gorm:"default:false" json:"is_core"`
	CitationCount         int            `gorm:"default:0" json:"citation_count"`
	Language              string         `gorm:"type:varchar(50)" json:"language"`
	JournalName           string         `gorm:"type:varchar(255)" json:"journal_name"`
	Volume                string         `gorm:"type:varchar(20)" json:"volume"`
	Issue                 string         `gorm:"type:varchar(20)" json:"issue"`
	StartPage             string         `gorm:"type:varchar(20)" json:"start_page"`
	EndPage               string         `gorm:"type:varchar(20)" json:"end_page"`
	AuthorType            string         `gorm:"type:varchar(50)" json:"author_type"`
	IsFirstAuthor         bool           `gorm:"default:false" json:"is_first_author"`
	IsCoFirstAuthor       bool           `gorm:"default:false" json:"is_co_first_author"`
	IsCorrespondingAuthor bool           `gorm:"default:false" json:"is_corresponding_author"`
	PublishDate           *time.Time     `json:"publish_date"`
	Status                string         `gorm:"type:varchar(20);default:'draft'" json:"status"`
	SubmitterID           uint           `gorm:"index;not null" json:"submitter_id"`
	SubmitTime            time.Time      `gorm:"not null" json:"submit_time"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Journal     *Journal      `gorm:"foreignKey:JournalID" json:"journal,omitempty"`
	Submitter   *User         `gorm:"foreignKey:SubmitterID" json:"submitter,omitempty"`
	Authors     []*Author     `gorm:"foreignKey:PaperID" json:"authors,omitempty"`
	Projects    []*Project    `gorm:"many2many:paper_projects;foreignKey:ID;references:ID" json:"projects,omitempty"`
	Attachments []*Attachment `gorm:"foreignKey:PaperID" json:"attachments,omitempty"`
	ReviewLogs  []*ReviewLog  `gorm:"foreignKey:PaperID" json:"review_logs,omitempty"`
}

// TableName 指定表名
func (Paper) TableName() string {
	return "papers"
}
