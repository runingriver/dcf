# 概述
一个DCF程序（自由现金流折现算法）：
1. 采用Golang实现DCF算法部分；
2. 自由现金流折现算法：
 - 输入参数：fcf 下一年自由现金流，单位：亿。总股本，单位：亿。贴现率：百分比。永续增长率：百分比。一阶段的年数N。未来N年的自由现金流年均增长率：百分比。
 - 输出：按照dcf计算过程输出每一步的结果：
  - 第一步：预测N年的自由现金流
  - 第二步：把这N年自由现金流折现成现值
  - 第三步：计算永续年金价值并把它折现成现值
  - 第四步：计算出企业价值，也就是N年的折现现金流加永续年金的折现值。
  - 第五步：计算每股价值
3. 支持多种交互方式：
 - 支持命令行交互；
 - 支持Web页面；
4. 系统架构设计：
 - 系统架构分两部分：前台展示，后台算法实现；
 - 请适当的加上注释，方便理解和维护；

# 整体说明
- 已在项目中实现完整的自由现金流折现（DCF）程序：后端使用 Golang 编写算法与交互，前端提供 Web 页面展示，亦支持命令行交互。
- 算法输出严格按照您要求的五步流程，包含每一步的数值结果与必要的校验（如 贴现率 > 永续增长率 ）。
- 代码已包含适当中文注释，便于理解与维护；同时附带基础单元测试验证关键逻辑。

# 使用方式
- 命令行交互：
  - 交互式输入： go run ./cmd/dcfcli
  - 使用参数直接计算： ./bin/dcfcli -fcf 10 -shares 100 -r 10 -gp 3 -n 5 -g 8
    - 参数含义： fcf 当前自由现金流（亿）， shares 总股本（亿）， r 贴现率（%）， gp 永续增长率（%）， n 一阶段年数， g 未来 N 年年均增长率（%）
- Web 页面：
  - 启动 Web 服务： go run ./cmd/server
  - 已启动本地服务：打开 http://localhost:8080/ 在浏览器输入参数并查看分步结果
  - 页面包含表单与详细计算过程，易于可视化分析
- 生成可执行文件
    - 生成 CLI 可执行文件： go build -o bin/dcfcli ./cmd/dcfcli
    - 生成 Web 服务可执行文件： go build -o bin/server ./cmd/server
    - 一次生成两者： go build -o bin/dcfcli ./cmd/dcfcli && go build -o bin/server ./cmd/server

# 算法说明

- 输入参数：
  - FCFBase : 当前自由现金流（亿）
  - TotalShares : 总股本（亿）
  - DiscountRatePct : 贴现率（%）
  - PerpetualGrowthPct : 永续增长率（%）
  - Years : 一阶段年数 N
  - AvgGrowthRatePct : 未来 N 年自由现金流年均增长率（%）
- 计算流程（均已实现且可视化输出）：
  - 第一步：预测 N 年自由现金流
    - FCF_t = FCF_{t-1} * (1 + g)
  - 第二步：折现 N 年自由现金流为现值
    - PV_t = FCF_t / (1 + r)^t
  - 第三步：计算永续年金价值并折现至现值
    - TV = FCF_N * (1 + gp) / (r - gp) ， PV_TV = TV / (1 + r)^N
  - 第四步：企业价值
    - FirmValue = Sum(PV_t) + PV_TV
  - 第五步：每股价值
    - PerShare = FirmValue / TotalShares
- 校验与边界：
  - 要求 r > gp ，否则返回错误
  - Years >= 1 ， TotalShares > 0

# 系统架构
- 分层设计：
  - 后台算法实现（Golang）： internal/dcf 包提供纯算法逻辑，独立且可复用
  - 前台展示：
    - 命令行应用： cmd/dcfcli 支持交互与参数两种模式
    - Web 服务： cmd/server 提供 HTTP 服务，模板渲染 web/index.html ，静态样式 web/style.css
- 依赖与实现：
  - 完全使用标准库（ net/http 、 html/template / text/template ），无需额外第三方依赖
  - 模板中提供函数 add 用于便捷地显示年序号

# 关键代码位置

- DCF 核心算法入口： internal/dcf/dcf.go:37 Compute(in Input) (StepResult, error)
- 参数结构体与结果结构体： internal/dcf/dcf.go:10 、 internal/dcf/dcf.go:21
- 预测与折现计算： internal/dcf/dcf.go:58 、 internal/dcf/dcf.go:66
- 永续年金与企业价值： internal/dcf/dcf.go:75 、 internal/dcf/dcf.go:79
- 每股价值： internal/dcf/dcf.go:82
- CLI 交互与输出： cmd/dcfcli/main.go:48 （交互）， cmd/dcfcli/main.go:95 （分步打印）
- Web 服务路由与模板： cmd/server/main.go:26 （路由）， cmd/server/main.go:19 （模板函数）， cmd/server/main.go:46 （计算与渲染）
- 前端模板与样式： web/index.html 、 web/style.css

# 验证与质量

- 单元测试： internal/dcf/dcf_test.go 覆盖基础正确性（长度、首年预测、每股价值一致性、异常条件）
- 构建与运行：
  - 运行测试： go test ./... （已通过）
  - 构建 CLI 与服务端： go build -o bin/dcfcli ./cmd/dcfcli && go build -o bin/server ./cmd/server
  - 启动服务： ./bin/server （已启动，预览地址如上）
- 晨星提供的验证：http://www.zhongguosou.com/zonghe/xianjinliu_guzhi.aspx