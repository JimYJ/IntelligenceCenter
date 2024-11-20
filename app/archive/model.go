package archive

type Archive struct {
	ID          int     `json:"id" db:"id"`                     // 主键
	ArchiveName string  `json:"archive_name" db:"archive_name"` // 档案名称
	FileCount   int     `json:"file_count" db:"file_count"`     // 档案文件数
	CreatedAt   string  `json:"created_at" db:"created_at"`     // 创建时间
	UpdatedAt   *string `json:"updated_at" db:"updated_at"`     // 更新时间，使用指针以支持 null
}

type ArchiveData struct {
	ArchiveName     string `json:"archive_name" db:"archive_name"` // 档案名称
	TaskCount       int    `json:"task_count" db:"-"`              // 关联任务总数
	ActiveTaskCount int    `json:"active_task_count" db:"-"`       // 关联活跃任务总数
	FileCount       int    `json:"file_count" db:"-"`              // 档案文件数
	// LLMSettingName  string `json:"llm_setting_name" db:"llm_setting_name"` // LLM设置名称
	// ApiType         uint8  `json:"api_type,string" db:"api_type"`          // API类型 1-智能小助手 Api 2-Ollama
	// ExtractionMode  uint8  `json:"extraction_mode" db:"extraction_mode"`   // 提取模式
	// ExtractionModel string `json:"extraction_model" db:"extraction_model"` // 提取模型

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
	TaskName          string  `json:"task_name" db:"task_name"`                   // 任务名称
	ArchiveName       string  `json:"archive_name" db:"archive_name"`             // 档案名称
	ResourceName      string  `json:"resource_num" db:"resource_num"`             // 资源数量
	TaskID            int     `json:"task_id" db:"task_id"`                       // 任务ID
	ArchiveID         int     `json:"archive_id" db:"archive_id"`                 // 所属的档案ID
	ExtractionMode    uint8   `json:"extraction_mode" db:"extraction_mode"`       // 提取模式
	ApiKeyID          int     `json:"api_key_id" db:"api_key_id"`                 // llm_api_settings 表ID
	ExtractionModel   string  `json:"extraction_model" db:"extraction_model"`     // 提取模型
	LLMSettingName    string  `json:"llm_setting_name" db:"llm_setting_name"`     // LLM设置名称
	ApiType           uint8   `json:"api_type" db:"api_type"`                     // API类型 1-OpenAI API Api 2-Ollama
	OriginContent     string  `json:"origin_content" db:"origin_content"`         // 文档原始内容
	ExtractionContent string  `json:"extraction_content" db:"extraction_content"` // 提取后内容
	IsExtracted       bool    `json:"is_extracted" db:"is_extracted"`             // 是否被提取 0否1是
	IsTranslated      bool    `json:"is_translated" db:"is_translated"`           // 是否被翻译 0否1是
	SrcURL            string  `json:"src_url" db:"src_url"`                       // 来源网址/来源文档地址
	CreatedAt         string  `json:"created_at" db:"created_at"`                 // 创建时间
	UpdatedAt         *string `json:"updated_at" db:"updated_at"`                 // 更新时间
}

type DocResource struct {
	ID             int64  `json:"id" db:"id"`
	DocID          int64  `json:"doc_id" db:"doc_id"`
	ArchiveID      int    `json:"archive_id" db:"archive_id"`
	ResourceType   int8   `json:"resource_type" db:"resource_type"`
	ResourcePath   string `json:"resource_path" db:"resource_path"`
	ResourceStatus uint8  `json:"resource_status" db:"resource_status"`
	ResourceSize   int    `json:"resource_size" db:"resource_size"`
	CreatedAt      string `json:"created_at" db:"created_at"`
}
