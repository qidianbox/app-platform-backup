package validator

import (
	"fmt"
)

// ValidateModuleConfig 验证模块配置
func ValidateModuleConfig(moduleCode string, config map[string]interface{}) error {
	if moduleCode == "" {
		return fmt.Errorf("module code is required")
	}
	// 简单验证，实际可以根据模块类型做更复杂的验证
	return nil
}
