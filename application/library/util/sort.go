package util

import "sort"

type SortOrder string

const (
	OrderAsc  = "asc"
	OrderDesc = "desc"
)

//用于针对[]map[string]interface{}类型的数据排序
//实现sort包Interface接口
type sortSliceMapStringInterface struct {
	Item  *[]map[string]interface{}
	Key   string
	Order SortOrder
}

func (ssmsi *sortSliceMapStringInterface) Len() int {
	return len(*ssmsi.Item)
}
func (ssmsi *sortSliceMapStringInterface) Less(i, j int) bool {
	if ssmsi.Order == OrderDesc {
		return Int((*ssmsi.Item)[i][ssmsi.Key]) > Int((*ssmsi.Item)[j][ssmsi.Key])
	}
	return Int((*ssmsi.Item)[i][ssmsi.Key]) < Int((*ssmsi.Item)[j][ssmsi.Key])
}
func (ssmsi *sortSliceMapStringInterface) Swap(i, j int) {
	(*ssmsi.Item)[i], (*ssmsi.Item)[j] = (*ssmsi.Item)[j], (*ssmsi.Item)[i]
}

func SortSliceMapStringInterface(item *[]map[string]interface{}, key string, order SortOrder) {
	ssmsi := &sortSliceMapStringInterface{item, key, order}
	sort.Sort(ssmsi)
}
