package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed index.html
var content embed.FS

type record struct {
	ID          int64   `json:"id"`
	Requirement string  `json:"requirement"`
	MR          string  `json:"mr"`
	MRURL       string  `json:"mrUrl"`
	Owner       string  `json:"owner"`
	Group       string  `json:"group"`
	Lines       int     `json:"lines"`
	Efficiency  int     `json:"efficiency"`
	Score       float64 `json:"score"`
	Problem     string  `json:"problem"`
	Date        string  `json:"date"`
	Status      string  `json:"status"`
}

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		http.ServeFileFS(w, r, content, "index.html")
	})

	mux.HandleFunc("/api/records", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		records, err := listRecords(r.Context(), db)
		if err != nil {
			http.Error(w, "failed to load records", http.StatusInternalServerError)
			log.Printf("list records: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(records); err != nil {
			log.Printf("encode records: %v", err)
		}
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           logRequest(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Printf("AI R&D dashboard listening on http://localhost:%s\n", port)
	log.Fatal(server.ListenAndServe())
}

func openDB() (*sql.DB, error) {
	path := os.Getenv("SQLITE_PATH")
	if path == "" {
		path = "dashboard.db"
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}

	if err := seed(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS records (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	requirement TEXT NOT NULL,
	mr TEXT NOT NULL,
	mr_url TEXT NOT NULL,
	owner TEXT NOT NULL,
	group_name TEXT NOT NULL,
	lines INTEGER NOT NULL,
	efficiency INTEGER NOT NULL,
	score REAL NOT NULL,
	problem TEXT NOT NULL,
	record_date TEXT NOT NULL,
	status TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_records_date ON records(record_date);
CREATE INDEX IF NOT EXISTS idx_records_owner ON records(owner);
CREATE INDEX IF NOT EXISTS idx_records_group ON records(group_name);
`)
	return err
}

func seed(db *sql.DB) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM records`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
INSERT INTO records (
	requirement, mr, mr_url, owner, group_name, lines, efficiency, score, problem, record_date, status
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range seedRecords {
		if _, err := stmt.Exec(
			item.Requirement,
			item.MR,
			item.MRURL,
			item.Owner,
			item.Group,
			item.Lines,
			item.Efficiency,
			item.Score,
			item.Problem,
			item.Date,
			item.Status,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func listRecords(ctx context.Context, db *sql.DB) ([]record, error) {
	rows, err := db.QueryContext(ctx, `
SELECT id, requirement, mr, mr_url, owner, group_name, lines, efficiency, score, problem, record_date, status
FROM records
ORDER BY record_date DESC, id DESC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []record
	for rows.Next() {
		var item record
		if err := rows.Scan(
			&item.ID,
			&item.Requirement,
			&item.MR,
			&item.MRURL,
			&item.Owner,
			&item.Group,
			&item.Lines,
			&item.Efficiency,
			&item.Score,
			&item.Problem,
			&item.Date,
			&item.Status,
		); err != nil {
			return nil, err
		}
		records = append(records, item)
	}

	return records, rows.Err()
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start).Round(time.Millisecond))
	})
}

var seedRecords = []record{
	{
		Requirement: "智能工单优先级排序",
		MR:          "feat: ranking pipeline with prompt guard",
		MRURL:       "https://git.example.com/ai-platform/workbench/-/merge_requests/1287",
		Owner:       "林越",
		Group:       "cos",
		Lines:       1860,
		Efficiency:  34,
		Score:       4.7,
		Problem:     "领域词召回不稳，补充了评测集与回退策略",
		Date:        "2026-07-08",
		Status:      "已合并",
	},
	{
		Requirement: "灰度发布配置生成器",
		MR:          "feat: rollout composer and policy preview",
		MRURL:       "https://git.example.com/ai-platform/workbench/-/merge_requests/1269",
		Owner:       "顾辰",
		Group:       "d5",
		Lines:       1320,
		Efficiency:  28,
		Score:       4.4,
		Problem:     "JSON Schema 边界提示过宽，增加模板约束",
		Date:        "2026-07-07",
		Status:      "已合并",
	},
	{
		Requirement: "接口回归用例自动补全",
		MR:          "test: generated contract coverage for payment APIs",
		MRURL:       "https://git.example.com/ai-platform/workbench/-/merge_requests/1261",
		Owner:       "陈珂",
		Group:       "cos",
		Lines:       2430,
		Efficiency:  41,
		Score:       4.8,
		Problem:     "Mock 数据命名混乱，人工统一语义标签",
		Date:        "2026-07-06",
		Status:      "验证中",
	},
	{
		Requirement: "研发知识库检索优化",
		MR:          "perf: hybrid search rerank and snippets",
		MRURL:       "https://git.example.com/ai-platform/workbench/-/merge_requests/1248",
		Owner:       "苏禾",
		Group:       "d5",
		Lines:       970,
		Efficiency:  22,
		Score:       4.1,
		Problem:     "历史文档格式不一，清洗耗时高于预期",
		Date:        "2026-07-04",
		Status:      "已合并",
	},
	{
		Requirement: "埋点指标口径校验",
		MR:          "feat: metric consistency assistant",
		MRURL:       "https://git.example.com/ai-platform/workbench/-/merge_requests/1239",
		Owner:       "邵宁",
		Group:       "cos",
		Lines:       760,
		Efficiency:  19,
		Score:       3.9,
		Problem:     "跨端字段别名太多，仍需产品侧确认口径",
		Date:        "2026-07-02",
		Status:      "已合并",
	},
	{
		Requirement: "前端异常聚类诊断",
		MR:          "feat: crash cluster insight panel",
		MRURL:       "https://git.example.com/ai-platform/workbench/-/merge_requests/1226",
		Owner:       "顾辰",
		Group:       "d5",
		Lines:       1510,
		Efficiency:  31,
		Score:       4.5,
		Problem:     "堆栈相似度阈值需要结合线上样本微调",
		Date:        "2026-06-30",
		Status:      "验证中",
	},
	{
		Requirement: "权限策略差异分析",
		MR:          "feat: policy diff explainer",
		MRURL:       "https://git.example.com/ai-platform/workbench/-/merge_requests/1212",
		Owner:       "林越",
		Group:       "cos",
		Lines:       1180,
		Efficiency:  26,
		Score:       4.3,
		Problem:     "策略 DSL 缺少注释，生成解释需人工复核",
		Date:        "2026-06-28",
		Status:      "已合并",
	},
	{
		Requirement: "构建失败自动诊断",
		MR:          "fix: ci failure classifier and hints",
		MRURL:       "https://git.example.com/ai-platform/workbench/-/merge_requests/1198",
		Owner:       "陈珂",
		Group:       "d5",
		Lines:       890,
		Efficiency:  24,
		Score:       4.2,
		Problem:     "少量日志噪声误判为依赖问题",
		Date:        "2026-06-25",
		Status:      "已合并",
	},
}
