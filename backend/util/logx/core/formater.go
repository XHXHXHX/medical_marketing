package core

import (
	"fmt"
	"net/url"
)

type Formater interface {
	Format(map[string]interface{}) map[string]interface{}
}

type RuleFormater struct {
}

func (rf *RuleFormater) useRule(source map[string]interface{}) map[string]interface{} {
	target := make(map[string]interface{})
	for key, val := range source {
		fType, ok := ruleDic[key]
		if false == ok {
			rf.appendExtData(&target, key, val)
			continue
		}
		switch fType {
		case FloatType:
			{
				res, err := tryFloat(val)
				if err == nil {
					target[key] = res
				} else {
					rf.appendExtData(&target, key, val)
				}
			}
		case DoubleType:
			{
				res, err := tryDouble(val)
				if err == nil {
					target[key] = res
				} else {
					rf.appendExtData(&target, key, val)
				}
			}
		case Int32Type:
			{
				res, err := tryInt32(val)
				if err == nil {
					target[key] = res
				} else {
					rf.appendExtData(&target, key, val)
				}
			}
		case Int64Type:
			{
				res, err := tryInt64(val)
				if err == nil {
					target[key] = res
				} else {
					rf.appendExtData(&target, key, val)
				}
			}
		case StringType:
			{
				res, err := tryString(val)
				// cut string type
				if err == nil {
					if len(res) > msgSize {
						res = res[:msgSize] + "..."
					}
					target[key] = res
				} else {
					rf.appendExtData(&target, key, val)
				}
			}
		case BoolType:
			{
				res, err := tryBool(val)
				if err == nil {
					target[key] = res
				} else {
					rf.appendExtData(&target, key, val)
				}
			}
		default:
			{
				rf.appendExtData(&target, key, val)
			}
		}
	}
	// cut ext data
	ext, ok := target["ext_data"]
	if ok {
		if len(ext.(string)) > extDataSize {
			target["ext_data"] = ext.(string)[:extDataSize] + "..."
		}
	}
	return target
}

func (rf *RuleFormater) Format(source map[string]interface{}) map[string]interface{} {
	return rf.useRule(FlatMap(source))
}

func (rf *RuleFormater) appendExtData(target *map[string]interface{}, key string, val interface{}) {
	pairs := fmt.Sprintf("%s=%s", url.PathEscape(key), url.PathEscape(fmt.Sprintf("%v", val)))
	ext, ok := (*target)["ext_data"]
	if ok {
		(*target)["ext_data"] = ext.(string) + "&" + pairs
	} else {
		(*target)["ext_data"] = pairs
	}
}
