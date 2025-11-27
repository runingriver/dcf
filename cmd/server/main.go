package main

import (
    "fmt"
    "log"
    "net/http"
    "path/filepath"
    "strconv"
    "text/template"

    "shares-dcf/internal/dcf"
)

var tmpl *template.Template

func main() {
    // Parse templates at startup.
    var err error
    tmpl, err = template.New("index").Funcs(template.FuncMap{
        "add": func(a, b int) int { return a + b },
    }).ParseFiles(filepath.Join("web", "index.html"))
    if err != nil {
        log.Fatalf("模板加载失败: %v", err)
    }

    http.HandleFunc("/", handleIndex)
    http.HandleFunc("/compute", handleCompute)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web"))))

    addr := ":8080"
    log.Printf("DCF Web服务已启动: http://localhost%s/", addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatal(err)
    }
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    data := map[string]any{"Result": nil, "Error": ""}
    _ = tmpl.Execute(w, data)
}

func handleCompute(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "参数解析失败", http.StatusBadRequest)
        return
    }

    in, err := parseInput(r)
    if err != nil {
        data := map[string]any{"Result": nil, "Error": err.Error()}
        _ = tmpl.Execute(w, data)
        return
    }

    res, err := dcf.Compute(in)
    if err != nil {
        data := map[string]any{"Result": nil, "Error": err.Error()}
        _ = tmpl.Execute(w, data)
        return
    }

    data := map[string]any{"Result": res, "Error": "", "Input": in}
    _ = tmpl.Execute(w, data)
}

func parseInput(r *http.Request) (dcf.Input, error) {
    getF := func(name string) (float64, error) {
        v := r.FormValue(name)
        if v == "" {
            return 0, fmt.Errorf("缺少参数: %s", name)
        }
        num, err := strconv.ParseFloat(v, 64)
        if err != nil {
            return 0, fmt.Errorf("参数格式错误: %s", name)
        }
        return num, nil
    }
    getI := func(name string) (int, error) {
        v := r.FormValue(name)
        if v == "" {
            return 0, fmt.Errorf("缺少参数: %s", name)
        }
        num, err := strconv.Atoi(v)
        if err != nil {
            return 0, fmt.Errorf("参数格式错误: %s", name)
        }
        return num, nil
    }

    fcf, err := getF("fcf")
    if err != nil { return dcf.Input{}, err }
    shares, err := getF("shares")
    if err != nil { return dcf.Input{}, err }
    rPct, err := getF("r")
    if err != nil { return dcf.Input{}, err }
    gpPct, err := getF("gp")
    if err != nil { return dcf.Input{}, err }
    years, err := getI("n")
    if err != nil { return dcf.Input{}, err }
    gPct, err := getF("g")
    if err != nil { return dcf.Input{}, err }

    return dcf.Input{
        FCFBase:            fcf,
        TotalShares:        shares,
        DiscountRatePct:    rPct,
        PerpetualGrowthPct: gpPct,
        Years:              years,
        AvgGrowthRatePct:   gPct,
    }, nil
}
