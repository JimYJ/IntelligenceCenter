package archive

type Archive struct {
	ID              int     `json:"id"`               // 主键
	ArchiveName     string  `json:"archive_name"`     // 档案名称
	FileCount       int     `json:"file_count"`       // 档案文件数
	ExtractionMode  uint8   `json:"extraction_mode"`  // 提取模式
	ApiKeyID        int     `json:"api_key_id"`       // llm_api_settings 表ID
	ExtractionModel string  `json:"extraction_model"` // 提取模型
	CreatedAt       string  `json:"created_at"`       // 创建时间
	UpdatedAt       *string `json:"updated_at"`       // 更新时间，使用指针以支持 null
}

type ArchiveDoc struct {
	ID                int     `json:"id"`                 // 主键
	DocName           string  `json:"doc_name"`           // 文档名称
	TaskID            int     `json:"task_id"`            // 任务ID
	ArchiveID         int     `json:"archive_id"`         // 所属的档案ID
	OriginContent     string  `json:"origin_content"`     // 文档原始内容
	ExtractionContent string  `json:"extraction_content"` // 提取后内容
	TranslateContent  string  `json:"translate_content"`  // 翻译后内容
	IsTranslated      bool    `json:"is_translated"`      // 是否被翻译 0否1是
	SrcURL            string  `json:"src_url"`            // 来源网址/来源文档地址
	CreatedAt         string  `json:"created_at"`         // 创建时间
	UpdatedAt         *string `json:"updated_at"`         // 更新时间
}
