package main

func defaultEqual(d1, d2 interface{}) bool {
	if d1 == nil && d2 == nil {
		return true
	}
	var bts1 []byte
	var bts2 []byte
	if d1 == nil {
		bts2 = d2.([]byte)
		if len(bts2) == 0 {
			return true
		}
		return false
	}
	if d2 == nil {
		bts1 = d1.([]byte)
		if len(bts1) == 0 {
			return true
		}
		return false
	}
	return string(bts1) == string(bts2)
}

func commentEqual(sComments, tComments map[string]string, name1, name2 string) bool {

	c1, ok1 := sComments[name1]
	c2, ok2 := tComments[name2]
	if !ok1 && !ok2 {
		return true
	}
	if !ok2 || !ok1 {
		return false
	}
	return c1 == c2
}

// 检查列是否存在
func columnExists(columns []*Column, columnName string) bool {
	for _, column := range columns {
		if column.Field == columnName {
			return true
		}
	}
	return false
}

func indexExists(indexes *[]Index, indexName string) bool {
	for _, idx := range *indexes {
		if idx.KeyName == indexName {
			return true
		}
	}
	return false
}

func tableExists(tables []string, table string) bool {
	for _, s := range tables {
		if s == table {
			return true
		}
	}
	return false
}

// 将字符串切片连接为一个字符串
func joinStrings(strings []string, separator string) string {
	if len(strings) == 0 {
		return ""
	}
	if len(strings) == 1 {
		return strings[0]
	}
	return strings[0] + separator + joinStrings(strings[1:], separator)
}
