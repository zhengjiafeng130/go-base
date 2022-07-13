package base

import "strings"

type node struct {
	pattern  string  // 待匹配路由,根节点
	part     string  // 路由的一段
	children []*node // 子节点
	isWild   bool    // 是否精确匹配, *表示匹配所有后续的路由, :表示匹配这一part的路由
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	var child *node
	if child = n.matchChild(part); child == nil {
		child = &node{
			part:     part,
			isWild:   isPartWild(part),
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height + 1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]

	for _, child := range  n.matchChildren(part) {
		result := child.search(parts, height + 1)
		if result != nil {
			return result
		}
	}
	return nil
}

func isPartWild(part string) bool {
	return part[0] == ':' || part[0] == '*'
}