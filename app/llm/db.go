package llm

import (
	"IntelligenceCenter/common/sqlite"
	"IntelligenceCenter/service/log"
	"fmt"
)

func create(r *Request) bool {
	_, err := sqlite.Conn().Exec("INSERT INTO llm_api_settings (name,api_type,api_url,api_key,timeout,request_rate_limit,remark) VALUES (?,?,?,?,?,?,?)",
		r.Name, r.ApiType, r.ApiURL, r.ApiKey, r.Timeout, r.RequestRateLimit, r.Remark)
	if err != nil {
		log.Info("新增LMM API设置出错:", err)
		return false
	}
	return true
}

func del(id string) bool {
	_, err := sqlite.Conn().Exec("DELETE FROM llm_api_settings WHERE id = ?;", id)
	if err != nil {
		log.Info("删除LMM API设置出错:", err)
		return false
	}
	return true
}

func edit(r *Request) bool {
	sql := `UPDATE llm_api_settings 
			SET name = ?, 
				api_type = ?, 
				api_url = ?, 
				api_key = ?, 
				timeout = ?, 
				request_rate_limit = ?, 
				remark = ? 
			WHERE id = ?;`
	_, err := sqlite.Conn().Exec(sql, r.Name, r.ApiType, r.ApiURL, r.ApiKey, r.Timeout, r.RequestRateLimit, r.Remark, r.ID)
	if err != nil {
		log.Info("编辑LMM API设置出错:", err)
		return false
	}
	return true
}

func listByPage(start, pageSize int, keyword string) []*Request {
	var searchSql string
	if len(keyword) != 0 {
		searchSql = "where name Like CONCAT('%',?,'%') or api_key Like CONCAT('%',?,'%') or api_url Like CONCAT('%',?,'%')"
	}
	sql := `SELECT
				id,
				name,
				api_type,
				api_url,
				api_key,
				timeout,
				request_rate_limit,
				remark
			FROM
				llm_api_settings 
				%s
				LIMIT ?,?;`
	sql = fmt.Sprintf(sql, searchSql)
	list := make([]*Request, 0)
	params := make([]any, 0)
	if len(keyword) != 0 {
		params = append(params, keyword, keyword, keyword)
	}
	params = append(params, start, pageSize)
	err := sqlite.Conn().Select(&list, sql, params...)
	if err != nil {
		log.Info("查询llm设置表出错:", err)
		return list
	}
	return list
}

// 获取记录总数
func countRecord(keyword string) int {
	var searchSql string
	if len(keyword) != 0 {
		searchSql = "where name Like CONCAT('%',?,'%') or api_key Like CONCAT('%',?,'%') or api_url Like CONCAT('%',?,'%')"
	}
	sql := `SELECT
				count(1)
			FROM
				llm_api_settings 
				%s;`
	sql = fmt.Sprintf(sql, searchSql)
	var num int
	params := make([]any, 0)
	if len(keyword) != 0 {
		params = append(params, keyword, keyword, keyword)
	}
	err := sqlite.Conn().Get(&num, sql, params...)
	if err != nil {
		log.Info("查询llm设置总数出错:", err)
		return num
	}
	return num
}

func list() []*Request {
	sql := `SELECT
				id,
				name,
				api_type
			FROM
				llm_api_settings;`
	list := make([]*Request, 0)
	err := sqlite.Conn().Select(&list, sql)
	if err != nil {
		log.Info("查询llm设置表出错:", err)
		return list
	}
	return list
}
