package util

//用于针对[]map[string]interface{}类型的数据排序
//实现sort包Interface接口
type SortSliceMapStringInterface struct {
	Item *[]map[string]interface{}
	Key  string
}

func (ssmsi *SortSliceMapStringInterface) Len() int {
	return len(*ssmsi.Item)
}
func (ssmsi *SortSliceMapStringInterface) Less(i, j int) bool {
	return Int((*ssmsi.Item)[i][ssmsi.Key]) > Int((*ssmsi.Item)[j][ssmsi.Key])
}
func (ssmsi *SortSliceMapStringInterface) Swap(i, j int) {
	(*ssmsi.Item)[i], (*ssmsi.Item)[j] = (*ssmsi.Item)[j], (*ssmsi.Item)[i]
}
