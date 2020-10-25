package load_balance

import (
	"errors"
	"strconv"
)

type WeightRoundRobinBalance struct {
	curIndex int
	rss      []*WeightNode
	rsw      []int
	//观察主体
	//conf LoadBalanceConf
}

type WeightNode struct {
	addr            string
	weight          int //权重值
	currentWeight   int //节点当前权重
	effectiveWeight int //有效权重
}

func (r *WeightRoundRobinBalance) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("param len need 2")
	}

	w, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return err
	}

	node := &WeightNode{
		addr: params[0],
		weight: int(w),
	}
	node.effectiveWeight = node.weight
	r.rss = append(r.rss, node)
	return nil
}

func (r *WeightRoundRobinBalance) Next() string {
	total := 0
	var best *WeightNode
	for i := 0; i < len(r.rss); i++ {
		n := r.rss[i]
		// 统计所有有效权重之和
		total += n.effectiveWeight
		// 变更节点临时权重为 临时权重+有效权重
		n.currentWeight += n.effectiveWeight
		// 有效权重默认与权重相同, 通讯异常时-1, 通讯成功+1
		// 直到恢复到weight大小
		if n.effectiveWeight < n.weight {
			n.effectiveWeight++
		}
		// 选择最大临时权重节点
		if best == nil || n.currentWeight > best.currentWeight {
			best = n
		}
	}

	if best == nil {
		return ""
	}
	best.currentWeight -= total
	return best.addr
}