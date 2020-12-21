package fancy

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string // 待匹配路由，例如 /p/:lang
	part     string // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool // 是否精确匹配，part 含有 : 或 * 时为true
}


func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}


// 第一个匹配成功的节点，用于插入树
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 用于搜索节点下的子节点
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
	/**
	   标记 parts 高度, 以便一层一层加入进树
	   递归出口条件 插入到最后一层 即 len(parts) == height
	 */
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	/** 某个路径 比如 api v1 */
	part := parts[height]

	/** 判断是否在树中，不在则插入到树 */
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	/** 递归调用 */
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {

	/** 递归出口条件 匹配到最后一层
	    即 len(parts) == height
	    或者路径中有*的 返回当前的路由节点 */
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	/**
	    获取路径去节点匹配，匹配到则返回该节点
	    示例: /hello/:name/login
	    第一层调用返回的树节点为 hello -> :name -> login
	    第二层调用返回的树节点为 :name -> login
	    第三层调用返回的树节点为 login
	    三次调用后达到递归调用条件返回 node
	*/
	part := parts[height]
	children := n.matchChildren(part)

	/** 示例匹配到 login 时, 我的输入为login1，
	    此时 children 为 nil
	    则跳出循环 return nil
	    进入到回溯阶段 看其它的节点是否有符合要求的
	    有返回 node 没有返回 nil 至此匹配结束
	*/
	for _, child := range children {
		/** DFS 深度优先 */
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}