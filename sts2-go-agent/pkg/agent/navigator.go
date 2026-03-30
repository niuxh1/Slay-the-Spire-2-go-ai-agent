package agent

import (
	"fmt"
	"math"

	"github.com/niuxh/sts2-go-agent/pkg/models"
)

// MapCost 定义了极端避战策略的节点权重
// 我们希望尽量走 ?, Rest, Shop，极力避开 Monster 和 Elite
func getNodeCost(nodeType string) int {
	switch nodeType {
	case "Monster":
		return 100 // 极高成本
	case "Elite":
		return 500 // 绝对避开
	case "Rest", "Shop":
		return 0   // 最优选择
	case "Event":
		return 10  // 次优选择
	case "Treasure", "Boss":
		return 0
	default:
		return 50
	}
}

// CalculateOptimalPath 使用简单的 BFS/Dijkstra 寻找到达地图顶端（Boss）的最小成本路径
func CalculateOptimalPath(mapState *models.MapState) string {
	if mapState == nil || len(mapState.Nodes) == 0 {
		return ""
	}

	// 1. 构建图和成本表
	// 简单起见，我们为当前 AvailableNodes 评估其子树的最小总成本
	
	nodeMap := make(map[string]models.FullNode)
	for _, n := range mapState.Nodes {
		key := fmt.Sprintf("%d-%d", n.Row, n.Col)
		nodeMap[key] = n
	}

	// 缓存中间结果 (Memoization)
	memo := make(map[string]int)

	var getMinCost func(row, col int) int
	getMinCost = func(row, col int) int {
		key := fmt.Sprintf("%d-%d", row, col)
		
		if val, exists := memo[key]; exists {
			return val
		}

		node, exists := nodeMap[key]
		if !exists {
			return math.MaxInt32 / 2
		}

		cost := getNodeCost(node.NodeType)

		if len(node.Children) == 0 {
			memo[key] = cost
			return cost
		}

		minChildCost := math.MaxInt32 / 2
		for _, child := range node.Children {
			childCost := getMinCost(child.Row, child.Col)
			if childCost < minChildCost {
				minChildCost = childCost
			}
		}

		total := cost + minChildCost
		memo[key] = total
		return total
	}

	bestOptionIndex := -1
	minCost := math.MaxInt32

	var analysis string
	for _, avail := range mapState.AvailableNodes {
		cost := getMinCost(avail.Row, avail.Col)
		analysis += fmt.Sprintf("Option %d (%s at %d,%d): Global Cost = %d\n", avail.Index, avail.NodeType, avail.Row, avail.Col, cost)
		if cost < minCost {
			minCost = cost
			bestOptionIndex = avail.Index
		}
	}

	if bestOptionIndex != -1 {
		return fmt.Sprintf("NAVIGATOR ALGORITHM:\n%s\n=> CONCLUSION: To achieve the EXTREME PACIFIST globally optimal path, you MUST choose option_index: %d.", analysis, bestOptionIndex)
	}

	return ""
}
