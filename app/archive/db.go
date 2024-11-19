package archive

import (
	"IntelligenceCenter/app/common"
	"IntelligenceCenter/common/sqlite"
	"IntelligenceCenter/service/log"
	"fmt"
)

func archiveListByPage(start, pageSize int, keyword string) []*Archive {
	var searchSql string
	if len(keyword) != 0 {
		searchSql = "where a.archive_name Like CONCAT('%',?,'%')"
	}
	sql := `SELECT 
				a.id,
				a.archive_name, 
				COUNT(ad.id) AS file_count,
				strftime('%s', a.created_at, '%s') created_at,
				strftime('%s', a.updated_at, '%s') updated_at
			FROM 
				archive a
			LEFT JOIN 
				archive_docs ad ON ad.archive_id = a.id
			%s
			GROUP BY
				a.id
			LIMIT ? ,?;`
	format := "%Y-%m-%d %H:%M:%S"
	sql = fmt.Sprintf(sql, format, common.GetTimeZone(), format, common.GetTimeZone(), searchSql)
	var list []*Archive
	params := make([]any, 0)
	if len(keyword) != 0 {
		params = append(params, keyword)
	}
	params = append(params, start, pageSize)
	err := sqlite.Conn().Select(&list, sql, params...)
	if err != nil {
		log.Info("查询档案表出错:", err)
		return list
	}
	return list
}

func archiveList() []*Archive {
	sql := `SELECT 
				a.id,
				a.archive_name
			FROM 
				archive a`
	var list []*Archive
	err := sqlite.Conn().Select(&list, sql)
	if err != nil {
		log.Info("查询档案表出错:", err)
		return list
	}
	return list
}

// 获取记录总数
func CountRecord(keyword string) int {
	var searchSql string
	if len(keyword) != 0 {
		searchSql = "where archive_name Like CONCAT('%',?,'%')"
	}
	sql := `SELECT
				count(1)
			FROM
				archive 
				%s;`
	sql = fmt.Sprintf(sql, searchSql)
	var num int
	params := make([]any, 0)
	if len(keyword) != 0 {
		params = append(params, keyword)
	}
	err := sqlite.Conn().Get(&num, sql, params...)
	if err != nil {
		log.Info("查询档案总数出错:", err)
		return num
	}
	return num
}

func docListByPage(start, pageSize int, id, keyword string) []*ArchiveDoc {
	var searchSql string
	if len(keyword) != 0 {
		searchSql = " AND doc_name Like CONCAT('%',?,'%')"
	}
	sql := `SELECT
				ad.id,
				ad.doc_name,
				t.task_name,
				a.archive_name,
                ad.is_extracted,
				ad.is_translated,
                count(dr.id) resource_num, 
				IFNULL(ad.extraction_mode,0) extraction_mode,
				IFNULL(ad.api_key_id,0) api_key_id,
				IFNULL(ad.extraction_model,0) extraction_model, 
                IFNULL(las.api_type,0) api_type,
				IFNULL(las.name,'') llm_setting_name,
				strftime('%s', ad.created_at, '%s') created_at,
				strftime('%s', ad.updated_at, '%s') updated_at
			FROM 
				archive_docs ad
			LEFT JOIN task t ON ad.task_id = t.id
			LEFT JOIN archive a ON ad.archive_id = a.id
            LEFT JOIN doc_resource dr ON dr.doc_id = ad.id
			LEFT JOIN llm_api_settings las ON las.id = ad.api_key_id 
			where ad.archive_id = ?
                %s
            group by ad.id
			LIMIT ? , ?;`
	format := "%Y-%m-%d %H:%M:%S"
	sql = fmt.Sprintf(sql, format, common.GetTimeZone(), format, common.GetTimeZone(), searchSql)
	var list []*ArchiveDoc
	err := sqlite.Conn().Select(&list, sql, id, start, pageSize)
	if err != nil {
		log.Info("查询文档表出错:", err)
		return list
	}
	return list
}

// 获取记录总数
func DocCountRecord(id, keyword string) int {
	var searchSql, idSql string
	if len(id) != 0 {
		idSql = " AND ad.archive_id = ? "
	}
	if len(keyword) != 0 {
		searchSql = " AND doc_name Like CONCAT('%',?,'%')"
	}
	sql := `SELECT
				count(1)
			FROM
				archive_docs ad
			where 1 = 1
			%s
            %s;`
	sql = fmt.Sprintf(sql, idSql, searchSql)
	var num int
	params := make([]any, 0)
	if len(id) != 0 {
		params = append(params, id)
	}
	if len(keyword) != 0 {
		params = append(params, keyword)
	}
	err := sqlite.Conn().Get(&num, sql, params...)
	if err != nil {
		log.Info("查询文档数出错:", err)
		return num
	}
	return num
}

// 档案相关信息
func archiveInfo(id string) *ArchiveData {
	sql := `SELECT
				archive_name
			FROM
				archive
			WHERE
				id = ?;`
	archiveData := &ArchiveData{}
	err := sqlite.Conn().Get(archiveData, sql, id)
	if err != nil {
		log.Info("查询档案信息出错:", err)
		return archiveData
	}
	return archiveData
}

// 获取记录总数
func archiveTask(id string, status int8) int {
	var statusSql string
	if status > 0 {
		statusSql = " AND task_status = ?"
	}
	sql := `SELECT
				count(1)
			FROM
				task 
			where archive_id = ?
				%s;`
	sql = fmt.Sprintf(sql, statusSql)
	var num int
	params := make([]any, 0)
	params = append(params, id)
	if status > 0 {
		params = append(params, status)
	}
	err := sqlite.Conn().Get(&num, sql, params...)
	if err != nil {
		log.Info("查询档案关联任务数出错:", err)
		return num
	}
	return num
}

// 创建
func Create(name string) int64 {
	sql := `INSERT INTO "archive" ( "archive_name" ) VALUES ( ? );`
	rows, err := sqlite.Conn().Exec(sql, name)
	if err != nil {
		log.Info("创建档案出错:", err)
		return -1
	}
	lastID, err := rows.LastInsertId()
	if err != nil {
		log.Info("获取新档案ID出错:", err)
		return -1
	}
	return lastID
}

// 获取记录总数
func DocResCountRecord(id string) int {
	var idSql string
	if len(id) != 0 {
		idSql = " AND dr.doc_id = ?"
	}
	sql := `SELECT
				count(1)
			FROM
				doc_resource dr
			where 1 = 1
			%s;`
	sql = fmt.Sprintf(sql, idSql)
	var num int
	params := make([]any, 0)
	if len(id) != 0 {
		params = append(params, id)
	}
	err := sqlite.Conn().Get(&num, sql, params...)
	if err != nil {
		log.Info("查询文档资源总数出错:", err)
		return num
	}
	return num
}

// 创建文档
func CreateDoc(taskID, archiveID int, name, content, srcUrl string) int64 {
	sql := `INSERT INTO "archive_docs" (
				"doc_name", 
				"task_id", 
				"archive_id", 
				"origin_content", 
				"is_extracted", 
				"is_translated", 
				"src_url" 
			) VALUES (
				?,  -- 文档名称
				?,  -- 任务ID
				?,  -- 档案ID
				?,  -- 原始内容
				0,  -- 是否被提取 (0或1)
				0,  -- 是否被翻译 (0或1)
				?   -- 来源网址
			);`
	rows, err := sqlite.Conn().Exec(sql, name, taskID, archiveID, content, srcUrl)
	if err != nil {
		log.Info("创建文档出错:", err)
		return -1
	}
	lastID, err := rows.LastInsertId()
	if err != nil {
		log.Info("获取新文档ID出错:", err)
		return -1
	}
	return lastID
}

// 批量创建文档
func CreateDocBatch(docs []*ArchiveDoc) error {
	sql := `INSERT INTO "archive_docs" (
				"doc_name", 
				"task_id", 
				"archive_id", 
				"origin_content", 
				"is_extracted", 
				"is_translated", 
				"src_url"
			) VALUES (
				:doc_name,
				:task_id,
				:archive_id, 
				:origin_content,
				:is_extracted,
				:is_translated,
				:src_url
			);`
	_, err := sqlite.Conn().NamedExec(sql, docs)
	if err != nil {
		log.Info("批量创建文档出错:", err)
		return err
	}
	return nil
}

// 更新文档
func UpdateDocByExtraction(docID, apiKeyID int, extractionMode uint8, extractionContent *string, extractionModel string) {
	sql := `UPDATE "archive_docs" SET
				"extraction_content" = ?,
				"extraction_mode" = ?,
				"api_key_id" = ?,
				"extraction_model" = ?,
				"is_extracted" = ?,
				"updated_at" = CURRENT_TIMESTAMP
			WHERE "id" = ?;`
	_, err := sqlite.Conn().Exec(sql, extractionContent, extractionMode, apiKeyID, extractionModel, 1, docID)
	if err != nil {
		log.Info("更新文档出错:", err)
	}
}

// 写入资源
func SaveDocResource(resources []*DocResource) error {
	sql := `INSERT INTO doc_resource (doc_id, resource_type, resource_path, resource_status, resource_size) 
			VALUES (:doc_id, :resource_type, :resource_path, :resource_status, :resource_size)`
	_, err := sqlite.Conn().NamedExec(sql, resources)
	if err != nil {
		log.Info("批量插入 doc_resource 表出错:", err)
		return err
	}
	return nil
}

// 获取资源
func GetDocResourceByDocID(docID int64) []*DocResource {
	sql := `SELECT
				id,
				doc_id,
				resource_type,
				resource_path,
				resource_status,
				resource_size,
				created_at
			FROM doc_resource
			WHERE doc_id = ?`
	var resources []*DocResource
	err := sqlite.Conn().Select(&resources, sql, docID)
	if err != nil {
		log.Info("查询 doc_resource 表出错:", err)
		return nil
	}
	return resources
}
