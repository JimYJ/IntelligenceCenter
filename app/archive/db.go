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
				a.extraction_mode,
				a.api_key_id,
				a.extraction_model, 
                las.api_type,
				las.name llm_setting_name,
				strftime('%s', a.created_at, '%s') created_at,
				strftime('%s', a.updated_at, '%s') updated_at
			FROM 
				archive a
			LEFT JOIN 
				archive_docs ad ON ad.archive_id = a.id
            LEFT JOIN llm_api_settings las ON las.id = a.api_key_id 
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

// 获取记录总数
func archiveCountRecord(keyword string) int {
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
				strftime('%s', ad.created_at, '%s') created_at,
				strftime('%s', ad.updated_at, '%s') updated_at
			FROM 
				archive_docs ad
			LEFT JOIN
				task t ON ad.task_id = t.id
			LEFT JOIN
				archive a ON ad.archive_id = a.id
            LEFT JOIN
				doc_resource dr ON dr.doc_id = ad.id
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
func docCountRecord(id, keyword string) int {
	var searchSql string
	if len(keyword) != 0 {
		searchSql = " AND doc_name Like CONCAT('%',?,'%')"
	}
	sql := `SELECT
				count(1)
			FROM
				archive_docs ad
			where ad.archive_id = ? 
            %s;`
	sql = fmt.Sprintf(sql, searchSql)
	var num int
	params := make([]any, 0)
	params = append(params, id)
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

// 档案相关信息
func archiveInfo(id string) *ArchiveData {
	sql := `SELECT
				a.extraction_mode,
				a.extraction_model,
				las.api_type,
				las.name llm_setting_name 
			FROM
				archive a
				LEFT JOIN llm_api_settings las ON las.id = a.api_key_id 
			WHERE
				a.id = ?
			GROUP BY
				a.id;`
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
