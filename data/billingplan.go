package data

// BillingFlags 用于设置计划授权使用的集成
type BillingFlags int

// BillingPlan 定义一个有权访问和设置限制的计划
type BillingPlan struct {
	ID          string                 `json:"id"`
	StripeID    string                 `json:"stripeId"`
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Price       float32                `json:"price"`
	YearlyPrice float32                `json:"yearly"`
	Params      map[string]interface{} `json:"params"`
}

var plans map[string]BillingPlan

func init() {
	plans = make(map[string]BillingPlan)
}

func AddPlan(plan BillingPlan) {
	plans[plan.ID] = plan
}

// GetPlan 通过 ID 获取计划
func GetPlan(id string) (BillingPlan, bool) {
	v, ok := plans[id]
	return v, ok
}

// GetPlans 返回所需版本计划的切片
func GetPlans(v string) []BillingPlan {
	var list []BillingPlan
	for k, p := range plans {
		if k == "free" {
			// 免费套餐适用于所有版本
			list = append(list, p)
		} else if p.Version == v {
			// 这是请求版本的计划
			list = append(list, p)
		}
	}
	return list
}

// GetPlansVersion 返回与当前计划匹配的计划的切片
func GetPlansVersion(plan string, defaultVersion string) []BillingPlan {
	if p, ok := plans[plan]; ok {
		return GetPlans(p.Version)
	}
	// 我们正在返回目前的计划，因为我们找不到这个计划
	return GetPlans(defaultVersion)
}
