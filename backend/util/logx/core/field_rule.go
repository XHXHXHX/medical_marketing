package core

import (
	"fmt"
	"sync"
)

var (
	ruleDic map[string]FieldType
	rLocker sync.Mutex
)

func init() {
	// 注册默认保留字段
	ruleDic = make(map[string]FieldType)
	ruleDic["x-user-id"] = StringType
	ruleDic["x-request-id"] = StringType
	ruleDic["x-trace-id"] = StringType
	ruleDic["project"] = StringType
}

func RegisterRule(key string, vType FieldType) error {
	rLocker.Lock()
	defer rLocker.Unlock()

	target, ok := ruleDic[key]
	if ok {
		return fmt.Errorf("Duplicate key:%s type:%d to %d", key, target, vType)
	}
	ruleDic[key] = vType
	return nil
}
