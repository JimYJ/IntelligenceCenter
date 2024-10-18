package llm

import (
	"IntelligenceCenter/common/sqlite"
	"log"
)

func save(r *Request) bool {
	_, err := sqlite.Conn().Exec("insert into llm_api_settings (name,api_type,api_url,api_key,timeout,request_rate_limit,remark) VALUE (?,?,?,?,?,?,?)",
		r.Name, r.ApiType, r.ApiURL, r.ApiKey, r.Timeout, r.RequestRateLimit, r.Remark)
	if err != nil {
		log.Println("新增LMM API设置出错:", err)
		return false
	}
	return true
}

func del(id string) bool {
	_, err := sqlite.Conn().Exec("DELETE FROM llm_api_settings WHERE id = ?;", id)
	if err != nil {
		log.Println("删除LMM API设置出错:", err)
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
		log.Println("编辑LMM API设置出错:", err)
		return false
	}
	return true
}
