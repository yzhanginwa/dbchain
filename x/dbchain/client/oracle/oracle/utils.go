package oracle

// UnWrap is mainly used to retrieve deep maps from result of function json.Unmarshal
func UnWrap(i interface{}, key string) interface{} {
   oneMap := i.(map[string]interface{})
   return oneMap[key]
}
