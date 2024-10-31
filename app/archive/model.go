package archive

type Archive struct {
	ID              int     `json:"id" db:"id"`                             // 主键
	ArchiveName     string  `json:"archive_name" db:"archive_name"`         // 档案名称
	FileCount       int     `json:"file_count" db:"file_count"`             // 档案文件数
	ExtractionMode  uint8   `json:"extraction_mode" db:"extraction_mode"`   // 提取模式
	ApiKeyID        int     `json:"api_key_id" db:"api_key_id"`             // llm_api_settings 表ID
	ExtractionModel string  `json:"extraction_model" db:"extraction_model"` // 提取模型
	CreatedAt       string  `json:"created_at" db:"created_at"`             // 创建时间
	UpdatedAt       *string `json:"updated_at" db:"updated_at"`             // 更新时间，使用指针以支持 null
}

type ArchiveData struct {
	ArchiveName     string `json:"archive_name" db:"archive_name"`         // 档案名称
	LLMSettingName  string `json:"name" db:"llm_setting_name"`             // LLM设置名称
	FileCount       int    `json:"file_count" db:"-"`                      // 档案文件数
	ExtractionMode  uint8  `json:"extraction_mode" db:"extraction_mode"`   // 提取模式
	ExtractionModel string `json:"extraction_model" db:"extraction_model"` // 提取模型
	TaskCount       int    `json:"task_count" db:"-"`                      // 关联任务总数
	ActiveTaskCount int    `json:"active_task_count" db:"-"`               // 关联活跃任务总数
}

type ArchiveTask struct {
	ArchiveName     string `json:"archive_name" db:"archive_name"`         // 档案名称
	LLMSettingName  string `json:"name" db:"llm_setting_name"`             // LLM设置名称
	FileCount       int    `json:"file_count" db:"-"`                      // 档案文件数
	ExtractionMode  uint8  `json:"extraction_mode" db:"extraction_mode"`   // 提取模式
	ExtractionModel string `json:"extraction_model" db:"extraction_model"` // 提取模型
	TaskCount       int    `json:"task_count" db:"-"`                      // 关联任务总数
	ActiveTaskCount int    `json:"active_task_count" db:"-"`               // 关联活跃任务总数
}

type ArchiveDoc struct {
	ID                int     `json:"id" db:"id"`                                 // 主键
	DocName           string  `json:"doc_name" db:"doc_name"`                     // 文档名称
	TaskID            int     `json:"task_id" db:"task_id"`                       // 任务ID
	ArchiveID         int     `json:"archive_id" db:"archive_id"`                 // 所属的档案ID
	OriginContent     string  `json:"origin_content" db:"origin_content"`         // 文档原始内容
	ExtractionContent string  `json:"extraction_content" db:"extraction_content"` // 提取后内容
	TranslateContent  string  `json:"translate_content" db:"translate_content"`   // 翻译后内容
	IsTranslated      bool    `json:"is_translated" db:"is_translated"`           // 是否被翻译 0否1是
	SrcURL            string  `json:"src_url" db:"src_url"`                       // 来源网址/来源文档地址
	CreatedAt         string  `json:"created_at" db:"created_at"`                 // 创建时间
	UpdatedAt         *string `json:"updated_at" db:"updated_at"`                 // 更新时间
}

type Keyword struct {
	Keyword string `db:"-" json:"keyword"` // 描述信息
}
