package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	proxy "github.com/shogo82148/go-sql-proxy"
	"log"
)

//func printArgs(args []driver.NamedValue) string {
//	sb := strings.Builder{}
//	sb.WriteString("[")
//	for idx, nv := range args {
//		if idx != 0 {
//			sb.WriteString(", ")
//		}
//		sb.WriteString(fmt.Sprintf("{%v $%d %v}", nv.Name, nv.Ordinal, nv.Value))
//	}
//	sb.WriteString("]")
//	return sb.String()
//}

func main() {
	//driverName := "sqlite3-proxy"
	//sql.Register(driverName, proxy.NewProxyContext(&sqlite3.SQLiteDriver{}, &proxy.HooksContext{
	//	PreExec: func(_ context.Context, _ *proxy.Stmt, _ []driver.NamedValue) (interface{}, error) {
	//		// The first return value(time.Now()) is passed to both `Hooks.Exec` and `Hook.ExecPost` callbacks.
	//		return time.Now(), nil
	//	},
	//	PostExec: func(_ context.Context, ctx interface{}, stmt *proxy.Stmt, args []driver.NamedValue, _ driver.Result, _ error) error {
	//		// The `ctx` parameter is the return value supplied from the `Hooks.PreExec` method, and may be nil.
	//		log.Printf(fmt.Sprintf("Exec query: %s; args = %v (%s)\n", stmt.QueryString, printArgs(args), time.Since(ctx.(time.Time))))
	//		return nil
	//	},
	//	PrePrepare: func(_ context.Context, _ *proxy.Stmt) (interface{}, error) {
	//		// The first return value(time.Now()) is passed to both `Hooks.Exec` and `Hook.ExecPost` callbacks.
	//		return time.Now(), nil
	//	},
	//	PostPrepare: func(_ context.Context, ctx interface{}, stmt *proxy.Stmt, _ error) error {
	//		// The `ctx` parameter is the return value supplied from the `Hooks.PreExec` method, and may be nil.
	//		log.Printf(fmt.Sprintf("Prepare query: %s (%s)\n", stmt.QueryString, time.Since(ctx.(time.Time))))
	//		return nil
	//	},
	//	PreQuery: func(_ context.Context, _ *proxy.Stmt, _ []driver.NamedValue) (interface{}, error) {
	//		// The first return value(time.Now()) is passed to both `Hooks.Exec` and `Hook.ExecPost` callbacks.
	//		return time.Now(), nil
	//	},
	//	PostQuery: func(_ context.Context, ctx interface{}, stmt *proxy.Stmt, args []driver.NamedValue, _ driver.Rows, _ error) error {
	//		// The `ctx` parameter is the return value supplied from the `Hooks.PreExec` method, and may be nil.
	//		log.Printf(fmt.Sprintf("Query query: %s; args = %v (%s)\n", stmt.QueryString, printArgs(args), time.Since(ctx.(time.Time))))
	//		return nil
	//	},
	//}))

	proxy.RegisterTracer()

	db, err := sql.Open("sqlite3:trace", ":memory:")
	if err != nil {
		log.Fatalf("Open filed: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS t1 (id INTEGER PRIMARY KEY)",
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(
		"INSERT INTO t1 (id) VALUES ($1)",
		42,
	)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(
		"INSERT INTO t1 (id) VALUES ($1)",
		44,
	)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT id FROM t1 WHERE id = ?1", 44)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	results := make([]int, 0)
	for rows.Next() {
		var i int
		err = rows.Scan(&i)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, i)
	}

	log.Printf("Result:")
	for _, ir := range results {
		log.Printf("row(%d)\n", ir)
	}
}
