package archive

import (
	"IntelligenceCenter/common/sqlite"
	"fmt"
	"log"
)

func archiveListByPage(start, pageSize int, keyword string) []*Archive {
	var searchSql string
	if len(keyword) != 0 {
		searchSql = "where archive_name Like CONCAT('%',?,'%')"
	}
	sql := `SELECT 
				a.id, 
				a.archive_name, 
				COUNT(ad.id) AS file_count,
				a.extraction_mode, 
				a.api_key_id, 
				a.extraction_model, 
				a.created_at, 
				a.updated_at
			FROM 
				archive a
			LEFT JOIN 
				archive_docs ad ON ad.archive_id = a.id
			%s
			GROUP BY 
				a.id 
			LIMIT ? ,?;`
	sql = fmt.Sprintf(sql, searchSql)
	list := make([]*Archive, 0)
	err := sqlite.Conn().Select(&list, sql, start, pageSize)
	if err != nil {
		log.Println("查询档案表出错:", err)
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
	err := sqlite.Conn().Select(&num, sql)
	if err != nil {
		log.Println("查询档案总数出错:", err)
		return num
	}
	return num
}

func docListByPage(start, pageSize int, keyword string) []*ArchiveDoc {
	var searchSql string
	if len(keyword) != 0 {
		searchSql = "where doc_name Like CONCAT('%',?,'%')"
	}
	sql := `SELECT 
				ad.id,
				ad.doc_name,
				ad.task_id, 
				t.task_name,
				ad.archive_id, 
				a.archive_name,
				ad.origin_content,
				ad.extraction_content,
				ad.translate_content,
				ad.is_translated,
				ad.src_url,
				ad.created_at,
				ad.updated_at
			FROM 
				archive_docs ad
			LEFT JOIN 
				task t ON ad.task_id = t.id
			LEFT JOIN 
				archive a ON ad.archive_id = a.id
			%s
			LIMIT ? , ?;`
	sql = fmt.Sprintf(sql, searchSql)
	list := make([]*ArchiveDoc, 0)
	err := sqlite.Conn().Select(&list, sql, start, pageSize)
	if err != nil {
		log.Println("查询档案表出错:", err)
		return list
	}
	return list
}

// 获取记录总数
func docCountRecord(keyword string) int {
	var searchSql string
	if len(keyword) != 0 {
		searchSql = "where doc_name Like CONCAT('%',?,'%')"
	}
	sql := `SELECT
				count(1)
			FROM
				archive_docs 
				%s;`
	sql = fmt.Sprintf(sql, searchSql)
	var num int
	err := sqlite.Conn().Select(&num, sql)
	if err != nil {
		log.Println("查询档案总数出错:", err)
		return num
	}
	return num
}
