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
				ad.extraction_mode,
				ad.api_key_id,
				ad.extraction_model, 
                las.api_type,
				las.name llm_setting_name,
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
