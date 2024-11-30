package core

///判断两个节点是否相等
func (this *ServiceNode) Equal(node *ServiceNode) bool {
	if node == nil {
		return false
	}
	if this.Tag != node.Tag {
		return false
	}
	if this.Type != node.Type {
		return false
	}
	if this.Id != node.Id {
		return false
	}
	if this.Version != node.Version {
		return false
	}
	if this.Addr != node.Addr {
		return false
	}
	if node.Meta != nil && this.Meta != nil && len(node.Meta) == len(this.Meta) {
		for k, v := range node.Meta {
			if v1, ok := this.Meta[k]; !ok || v != v1 {
				return false
			}
		}
	} else if node.Meta == nil && this.Meta == nil {
		return true
	} else {
		return false
	}
	return true
}
