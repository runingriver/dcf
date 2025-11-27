package dcf

import (
    "errors"
    "math"
)

// Input represents the parameters required to run a DCF valuation.
// 所有百分比输入均以百分数形式传入，例如10表示10%。
type Input struct {
    FCFBase             float64 // 下一年自由现金流（亿）
    TotalShares         float64 // 总股本（亿）
    DiscountRatePct     float64 // 贴现率（%）
    PerpetualGrowthPct  float64 // 永续增长率（%）
    Years               int     // 一阶段年数 N
    AvgGrowthRatePct    float64 // 未来N年自由现金流年均增长率（%）
}

// StepResult captures each step of the DCF computation for detailed output.
// 计算过程的每一步结果。
type StepResult struct {
    ProjectedFCF       []float64 // 预测的N年自由现金流（每年，单位：亿）
    DiscountedFCF      []float64 // 折现后的N年自由现金流现值（单位：亿）
    TerminalValue      float64   // 永续年金价值（发生在第N年末，单位：亿）
    DiscountedTerminal float64   // 永续年金折现至当期的现值（单位：亿）
    FirmValue          float64   // 企业价值（单位：亿）
    PerShareValue      float64   // 每股价值（单位：元/股，若股本单位为亿股）
}

// Compute performs the DCF valuation and returns step-by-step results.
// 算法遵循：
// 1) FCF_t = FCF_{t-1} * (1 + g)；
// 2) PV_t = FCF_t / (1 + r)^t；
// 3) TV_N = FCF_N * (1 + g_p) / (r - g_p)；PV_TV = TV_N / (1 + r)^N；
// 4) Firm = Sum(PV_t) + PV_TV；
// 5) PerShare = Firm / TotalShares。
func Compute(in Input) (StepResult, error) {
    if in.Years <= 0 {
        return StepResult{}, errors.New("years must be >= 1")
    }
    if in.TotalShares <= 0 {
        return StepResult{}, errors.New("total shares must be > 0")
    }
    // Convert percentage inputs to decimals.
    r := in.DiscountRatePct / 100.0
    g := in.AvgGrowthRatePct / 100.0
    gp := in.PerpetualGrowthPct / 100.0

    if r <= gp {
        return StepResult{}, errors.New("discount rate must be greater than perpetual growth rate")
    }

    projected := make([]float64, in.Years)
    discounted := make([]float64, in.Years)

    // Step 1: forecast N years of FCF.
    // Interpret FCFBase as next year's FCF (t=1).
    current := in.FCFBase
    for i := 0; i < in.Years; i++ {
        projected[i] = current
        current = current * (1.0 + g)
    }

    // Step 2: discount each year's FCF to present value.
    var sumPV float64
    for t := 1; t <= in.Years; t++ {
        ft := projected[t-1]
        pv := ft / math.Pow(1.0+r, float64(t))
        discounted[t-1] = pv
        sumPV += pv
    }

    // Step 3: terminal value at end of year N and discount to present.
    fcfN := projected[in.Years-1]
    terminal := fcfN * (1.0 + gp) / (r - gp)
    discountedTerminal := terminal / math.Pow(1.0+r, float64(in.Years))

    // Step 4: firm value.
    firm := sumPV + discountedTerminal

    // Step 5: per-share value.
    perShare := firm / in.TotalShares

    return StepResult{
        ProjectedFCF:       projected,
        DiscountedFCF:      discounted,
        TerminalValue:      terminal,
        DiscountedTerminal: discountedTerminal,
        FirmValue:          firm,
        PerShareValue:      perShare,
    }, nil
}
