package dcf

import "testing"

func TestComputeBasic(t *testing.T) {
    in := Input{
        FCFBase:            10,   // 亿
        TotalShares:        100,  // 亿股
        DiscountRatePct:    10,   // %
        PerpetualGrowthPct: 3,    // %
        Years:              5,
        AvgGrowthRatePct:   8,    // %
    }
    res, err := Compute(in)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(res.ProjectedFCF) != in.Years {
        t.Fatalf("projected length = %d, want %d", len(res.ProjectedFCF), in.Years)
    }
    if len(res.DiscountedFCF) != in.Years {
        t.Fatalf("discounted length = %d, want %d", len(res.DiscountedFCF), in.Years)
    }
    // First year projection equals input FCF (next year's FCF).
    wantYear1 := in.FCFBase
    if diff := abs(res.ProjectedFCF[0]-wantYear1); diff > 1e-9 {
        t.Fatalf("year1 projected diff=%g", diff)
    }
    // Per-share should equal firm value divided by total shares.
    if diff := abs(res.PerShareValue - res.FirmValue/in.TotalShares); diff > 1e-9 {
        t.Fatalf("per-share mismatch diff=%g", diff)
    }
}

func TestInvalidRates(t *testing.T) {
    in := Input{FCFBase: 10, TotalShares: 10, DiscountRatePct: 5, PerpetualGrowthPct: 5, Years: 3, AvgGrowthRatePct: 5}
    if _, err := Compute(in); err == nil {
        t.Fatalf("expected error when r <= gp")
    }
}

func abs(x float64) float64 { if x < 0 { return -x }; return x }

func TestCaseProvided(t *testing.T) {
    in := Input{
        FCFBase:            988.0000,
        TotalShares:        12.5200,
        DiscountRatePct:    10.00,
        PerpetualGrowthPct: 3.00,
        Years:              10,
        AvgGrowthRatePct:   8.00,
    }
    res, err := Compute(in)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    wantFirm := 19485.7206
    wantPerShare := 1556.367457
    if diff := abs(res.FirmValue - wantFirm); diff > 1e-3 {
        t.Fatalf("firm value diff=%g, got=%.6f want=%.6f", diff, res.FirmValue, wantFirm)
    }
    if diff := abs(res.PerShareValue - wantPerShare); diff > 1e-6 {
        t.Fatalf("per-share diff=%g, got=%.6f want=%.6f", diff, res.PerShareValue, wantPerShare)
    }
}
