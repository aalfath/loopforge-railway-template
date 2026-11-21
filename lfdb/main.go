// lfdb is the small database helper of the LoopForge Railway template.
//
// LoopForge creates its tables itself (gorm AutoMigrate), but its starter
// catalog, welcome project and two generated columns are plain SQL scripts that
// upstream's installer applies with the mysql client. The platform image has no
// mysql client, so this helper applies the same scripts, unchanged, over
// AGENTFLOW_MYSQL_DSN.
//
//	lfdb wait                  wait until MySQL answers (up to 5 minutes)
//	lfdb accounts              print the number of local password accounts
//	lfdb extensions SQL_DIR    add the missing generated columns, then apply the seeds
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

func open() (*sql.DB, error) {
	dsn := strings.TrimSpace(os.Getenv("AGENTFLOW_MYSQL_DSN"))
	if dsn == "" {
		return nil, errors.New("AGENTFLOW_MYSQL_DSN is not set")
	}
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("AGENTFLOW_MYSQL_DSN: %w", err)
	}
	cfg.MultiStatements = true
	cfg.ParseTime = true
	if cfg.Params == nil {
		cfg.Params = map[string]string{}
	}
	cfg.Params["charset"] = "utf8mb4"
	cfg.Timeout = 5 * time.Second
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

func wait(db *sql.DB) error {
	var last error
	for i := 0; i < 150; i++ {
		if last = db.Ping(); last == nil {
			return nil
		}
		if i%10 == 0 {
			fmt.Fprintf(os.Stderr, "lfdb: waiting for MySQL: %v\n", last)
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("MySQL did not answer within 5 minutes: %w", last)
}

func accounts(db *sql.DB) (int64, error) {
	var n int64
	err := db.QueryRow("SELECT COUNT(*) FROM local_accounts").Scan(&n)
	var me *mysql.MySQLError
	if errors.As(err, &me) && me.Number == 1146 { // table does not exist yet
		return 0, nil
	}
	return n, err
}

func extensions(db *sql.DB, dir string) error {
	obs, err := os.ReadFile(filepath.Join(dir, "001_run_obs_generated_columns.sql"))
	if err != nil {
		return err
	}
	for _, tc := range [][2]string{{"run_telemetry", "adapter"}, {"run_logs", "level"}} {
		var n int
		if err := db.QueryRow("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name=? AND column_name=?", tc[0], tc[1]).Scan(&n); err != nil {
			return err
		}
		if n == 1 {
			fmt.Printf("lfdb: already installed: %s.%s\n", tc[0], tc[1])
			continue
		}
		stmt := regexp.MustCompile(`(?s)ALTER TABLE ` + tc[0] + `\b.*?;`).Find(obs)
		if stmt == nil {
			return fmt.Errorf("missing generated-column definition for %s", tc[0])
		}
		if _, err := db.Exec(string(stmt)); err != nil {
			return fmt.Errorf("%s.%s: %w", tc[0], tc[1], err)
		}
		fmt.Printf("lfdb: installed: %s.%s\n", tc[0], tc[1])
	}
	// Same order as upstream's install-schema-extensions.sh: seeds, then 031.
	seeds, _ := filepath.Glob(filepath.Join(dir, "seed_*.sql"))
	sort.Strings(seeds)
	for _, f := range append(seeds, filepath.Join(dir, "031_generic_backend_openssh.sql")) {
		body, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(body)); err != nil {
			return fmt.Errorf("%s: %w", filepath.Base(f), err)
		}
		fmt.Printf("lfdb: applied %s\n", filepath.Base(f))
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: lfdb wait | accounts | extensions SQL_DIR")
		os.Exit(2)
	}
	db, err := open()
	if err == nil {
		defer db.Close()
		switch os.Args[1] {
		case "wait":
			err = wait(db)
		case "accounts":
			var n int64
			if n, err = accounts(db); err == nil {
				fmt.Println(n)
			}
		case "extensions":
			if len(os.Args) < 3 {
				err = errors.New("usage: lfdb extensions SQL_DIR")
			} else {
				err = extensions(db, os.Args[2])
			}
		default:
			err = fmt.Errorf("unknown command %q", os.Args[1])
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "lfdb:", err)
		os.Exit(1)
	}
}
