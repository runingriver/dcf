package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strconv"
    "strings"

    "shares-dcf/internal/dcf"
)

func main() {
    // Flags allow non-interactive usage.
    fcf := flag.Float64("fcf", 0, "下一年自由现金流（亿）")
    shares := flag.Float64("shares", 0, "总股本（亿）")
    r := flag.Float64("r", 0, "贴现率（%）")
    gp := flag.Float64("gp", 0, "永续增长率（%）")
    years := flag.Int("n", 0, "一阶段年数 N")
    g := flag.Float64("g", 0, "未来N年自由现金流年均增长率（%）")
    flag.Parse()

    var in dcf.Input
    // If all flags provided, use them; otherwise fall back to interactive.
    if *fcf > 0 && *shares > 0 && *r > 0 && *years > 0 {
        in = dcf.Input{
            FCFBase:            *fcf,
            TotalShares:        *shares,
            DiscountRatePct:    *r,
            PerpetualGrowthPct: *gp,
            Years:              *years,
            AvgGrowthRatePct:   *g,
        }
    } else {
        in = promptInteractive()
    }

    result, err := dcf.Compute(in)
    if err != nil {
        fmt.Fprintf(os.Stderr, "错误: %v\n", err)
        os.Exit(1)
    }

    printSteps(in, result)
}

func promptInteractive() dcf.Input {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("请输入以下参数，单位见括号提示（直接输入数值即可）：")

    fcf := promptFloat(reader, "下一年自由现金流（亿）: ")
    shares := promptFloat(reader, "总股本（亿）: ")
    r := promptFloat(reader, "贴现率（%）: ")
    gp := promptFloat(reader, "永续增长率（%）: ")
    years := promptInt(reader, "一阶段年数 N: ")
    g := promptFloat(reader, "未来N年自由现金流年均增长率（%）: ")

    return dcf.Input{
        FCFBase:            fcf,
        TotalShares:        shares,
        DiscountRatePct:    r,
        PerpetualGrowthPct: gp,
        Years:              years,
        AvgGrowthRatePct:   g,
    }
}

func promptFloat(r *bufio.Reader, label string) float64 {
    for {
        fmt.Print(label)
        text, _ := r.ReadString('\n')
        text = strings.TrimSpace(text)
        v, err := strconv.ParseFloat(text, 64)
        if err == nil {
            return v
        }
        fmt.Println("请输入合法的数值")
    }
}

func promptInt(r *bufio.Reader, label string) int {
    for {
        fmt.Print(label)
        text, _ := r.ReadString('\n')
        text = strings.TrimSpace(text)
        v, err := strconv.Atoi(text)
        if err == nil {
            return v
        }
        fmt.Println("请输入合法的整数")
    }
}

func printSteps(in dcf.Input, res dcf.StepResult) {
    fmt.Println()
    fmt.Println("按照DCF计算过程输出：")
    fmt.Printf("第一步：预测N年的自由现金流 (单位：亿)\n")
    for i, v := range res.ProjectedFCF {
        fmt.Printf("  第%d年：%.4f\n", i+1, v)
    }

    fmt.Printf("第二步：把这N年自由现金流折现成现值 (单位：亿)\n")
    for i, v := range res.DiscountedFCF {
        fmt.Printf("  第%d年现值：%.4f\n", i+1, v)
    }

    fmt.Printf("第三步：计算永续年金价值并把它折现成现值\n")
    fmt.Printf("  永续年金价值（第N年末）：%.4f (亿)\n", res.TerminalValue)
    fmt.Printf("  折现到当期的现值：%.4f (亿)\n", res.DiscountedTerminal)

    fmt.Printf("第四步：计算企业价值\n")
    fmt.Printf("  企业价值：%.4f (亿)\n", res.FirmValue)

    fmt.Printf("第五步：计算每股价值\n")
    fmt.Printf("  每股价值：%.6f\n", res.PerShareValue)

    fmt.Println()
    fmt.Printf("输入摘要：下一年FCF=%.4f亿, 股本=%.4f亿股, r=%.2f%%, g=%.2f%%, gp=%.2f%%, N=%d\n",
        in.FCFBase, in.TotalShares, in.DiscountRatePct, in.AvgGrowthRatePct, in.PerpetualGrowthPct, in.Years)
}
