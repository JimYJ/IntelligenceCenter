package task

import (
	"IntelligenceCenter/common/utils"
	"IntelligenceCenter/service/log"
	"os"

	"gopkg.in/yaml.v2"
)

func Match() {

}

func getRules() ([]*extractionRules, []string) {
	list, err := utils.FindFilesBySuffix("./extraction-rules", "yml", "yaml")
	log.Info(list, err)
	rulesList := make([]*extractionRules, 0)
	domainList := make([]string, 0)
	for _, item := range list {
		rules := getRulesFiles(item)
		if rules != nil {
			rulesList = append(rulesList, rules)
			domainList = append(domainList, rules.MatchDomain...)
		}
	}
	return rulesList, domainList
}

func getRulesFiles(path string) *extractionRules {
	rules := &extractionRules{}
	file, err := os.ReadFile(path)
	if err != nil {
		log.Info("读取文件错误:", err)
		return nil
	}
	err = yaml.Unmarshal(file, rules)
	if err != nil {
		log.Info("解析文件错误:", err)
		return nil
	}
	return rules
}
