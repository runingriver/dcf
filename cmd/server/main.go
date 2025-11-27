package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"strconv"

	"dcf/internal/dcf"
)

var tmpl *template.Template

func main() {
	// Parse templates at startup.
	var err error
	tmpl, err = template.New("index.html").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		// 提供 pow 以便在模板中演示公式中的幂次
		"pow": func(base float64, exp int) float64 { return math.Pow(base, float64(exp)) },
		// 计算 1 + p/100，用于 (1+r)^t 的底数
		"onePlusPct": func(p float64) float64 { return 1.0 + p/100.0 },
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := map[string]any{"Result": nil, "Error": ""}
	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, fmt.Sprintf("模板渲染失败: %v", err), http.StatusInternalServerError)
		return
	}
}

func handleCompute(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "参数解析失败", http.StatusBadRequest)
		return
	}

	in, err := parseInput(r)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// Preserve user input values in the form when parse error occurs.
		toF := func(name string) float64 {
			v := r.FormValue(name)
			num, _ := strconv.ParseFloat(v, 64)
			return num
		}
		toI := func(name string) int {
			v := r.FormValue(name)
			num, _ := strconv.Atoi(v)
			return num
		}
		preserved := dcf.Input{
			FCFBase:            toF("fcf"),
			TotalShares:        toF("shares"),
			DiscountRatePct:    toF("r"),
			PerpetualGrowthPct: toF("gp"),
			Years:              toI("n"),
			AvgGrowthRatePct:   toF("g"),
		}
		data := map[string]any{"Result": nil, "Error": err.Error(), "Input": preserved}
		_ = tmpl.ExecuteTemplate(w, "index.html", data)
		return
	}

	res, err := dcf.Compute(in)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		data := map[string]any{"Result": nil, "Error": err.Error(), "Input": in}
		_ = tmpl.ExecuteTemplate(w, "index.html", data)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := map[string]any{"Result": res, "Error": "", "Input": in}
	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, fmt.Sprintf("模板渲染失败: %v", err), http.StatusInternalServerError)
		return
	}
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
	if err != nil {
		return dcf.Input{}, err
	}
	shares, err := getF("shares")
	if err != nil {
		return dcf.Input{}, err
	}
	rPct, err := getF("r")
	if err != nil {
		return dcf.Input{}, err
	}
	gpPct, err := getF("gp")
	if err != nil {
		return dcf.Input{}, err
	}
	years, err := getI("n")
	if err != nil {
		return dcf.Input{}, err
	}
	gPct, err := getF("g")
	if err != nil {
		return dcf.Input{}, err
	}

	return dcf.Input{
		FCFBase:            fcf,
		TotalShares:        shares,
		DiscountRatePct:    rPct,
		PerpetualGrowthPct: gpPct,
		Years:              years,
		AvgGrowthRatePct:   gPct,
	}, nil
}
