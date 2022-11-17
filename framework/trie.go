package framework

import (
	"errors"
	"strings"
)

//树结构
type Tree struct {
	root *node //根节点
}

// 路由节点
type node struct {
	isLast   bool                //代表这个节点是否可以成为路由
	segment  string              //路由规格中某个节点		例如/user/login路由里面的 login
	handlers []ControllerHandler // 执行业务逻辑的ControllerHandler 和 中间件 ControllerHandler
	childs   []*node             //子节点

	parent *node //父节点
}

func newNode() *node {
	return &node{
		isLast:  false,
		segment: "",
		childs:  []*node{},
	}
}

func NewTree() *Tree {
	root := newNode()
	return &Tree{root: root}
}

// 判断是否是通配符 例如 ":id"
func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

//过滤所有匹配的子节点
func (n *node) filterChildNodes(segment string) []*node {
	if len(n.childs) == 0 {
		//如果没有子节点 返回nil
		return nil
	}

	if isWildSegment(segment) { //如果是通配符 则所有的子节点都匹配
		return n.childs
	}

	nodes := make([]*node, 0, len(n.childs)-1)

	for _, cnode := range n.childs {
		if isWildSegment(cnode.segment) { //如果子节点是通配符 则匹配
			nodes = append(nodes, cnode)
		} else if cnode.segment == segment { //如果不是通配符 但是段完全相同 则匹配
			nodes = append(nodes, cnode)
		}
	}

	return nodes
}

// 匹配节点
func (n *node) matchNode(url string) *node {
	// 将url分成两个段 取第一个段向下匹配
	segments := strings.SplitN(url, "/", 2)
	//取第一个段
	segment := segments[0]
	if !isWildSegment(segment) {
		// 如果不是通配符 将其转换为大写的
		segment = strings.ToUpper(segment)
	}

	cnodes := n.filterChildNodes(segment)
	if cnodes == nil || len(cnodes) == 0 {
		return nil
	}

	// 如果只有一段 则需是终极节点
	if len(segments) == 1 {
		for _, ln := range cnodes {
			if ln.isLast {
				return ln
			}
		}

		// 都不是最后一个节点 则不匹配
		return nil
	}

	// 两个以上的segment
	for _, ln := range cnodes {
		lnMatch := ln.matchNode(segments[1])
		if lnMatch != nil {
			return lnMatch
		}
	}

	return nil
}

func (tree *Tree) AddRouter(uri string, handlers []ControllerHandler) error {
	n := tree.root
	if n.matchNode(uri) != nil {
		return errors.New("router uri :" + uri + " is exsit")
	}

	segments := strings.Split(uri, "/")

	for index, segment := range segments {
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}

		var objNode *node //满足条件的下级节点

		var isLast = index == len(segments)-1

		childNodes := n.filterChildNodes(segment)
		if len(childNodes) > 0 {
			for _, cnode := range childNodes {
				if cnode.segment == segment {
					objNode = cnode
					break
				}
			}
		}

		if objNode == nil {
			newNode := newNode()
			newNode.segment = segment

			if isLast {
				newNode.isLast = true
				newNode.handlers = handlers
			}

			newNode.parent = n

			n.childs = append(n.childs, newNode)
			objNode = newNode
		}

		n = objNode
	}
	return nil
}

// 找控制器
func (tree *Tree) FindControllerHandler(uri string) []ControllerHandler {
	matchNode := tree.root.matchNode(uri)
	if matchNode == nil {
		return nil
	}

	return matchNode.handlers
}

// 最后节点 uri和完整的路由匹配，找出params参数
// 例如：/user/:id  /user/1  返回map["id"]{1}
func (n *node) parseParamsFromEndNode(uri string) map[string]string {
	segments := strings.Split(uri, "/")
	cur := n
	ret := make(map[string]string)

	for i := len(segments) - 1; i >= 0; i-- {
		if cur.segment == "" { //如果到了root节点 跳出循环
			break
		}
		if isWildSegment(cur.segment) { //通配符就赋值  例如：["a"]=1
			ret[cur.segment[1:]] = segments[i]
		}
		cur = cur.parent //查找父级节点
	}

	return ret
}
