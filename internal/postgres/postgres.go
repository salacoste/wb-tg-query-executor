package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/shomali11/xsql"
)

type Task struct {
	Id uint64
	BotToken string
	ChatId int64
	ChatDescribe string
	SqlQuery string
	ScheduleCron string
	LastExecutionTs sql.NullTime
}

func GetTaskList(db *sql.DB) ([]Task, error) {
	q := `
		SELECT id, bot_token, bot_chat_id, chat_describe, sql_query, schedule_cron, last_execution_ts
		FROM tg_query_executor
	`
	rows, err := db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("processing query: %s. error: %s", q, err)
	}
	defer rows.Close()
	
	tasks := []Task{}
	for rows.Next() {
		t := Task{}
		err := rows.Scan(&t.Id, &t.BotToken, &t.ChatId, &t.ChatDescribe, &t.SqlQuery, &t.ScheduleCron, &t.LastExecutionTs)
		if err != nil {
			return nil, fmt.Errorf("scan row: %s", err)
		}
		fixTimezone(&t.LastExecutionTs)
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func UpdateLastExecutionTs(db *sql.DB, t Task) error {
	q := `
		UPDATE tg_query_executor
		SET last_execution_ts = $2
		WHERE id = $1
	`
	_, err := db.Exec(q, t.Id, t.LastExecutionTs)
	if err != nil {
		return fmt.Errorf("processing query: %s. error: %s", q, err)
	}
	return nil
}

func ExecQuery(db *sql.DB, q string) (string, error) {
	rows, err := db.Query(q)
	if err != nil {
		return "", fmt.Errorf("query: %s error: %s", q, err)
	}
	defer rows.Close()
	return xsql.Pretty(rows)

	// res := QueryResult{}
	// res.Headers, err = rows.Columns()
	// log.Println(res.Headers)
	// if err != nil {
	// 	return QueryResult{}, fmt.Errorf("could not get columns from result of query: '%s' error: %s", q, err)
	// }

	// values := make([]sql.RawBytes, len(res.Headers))
	// valuePtrs := make([]interface{}, len(res.Headers))
	// for i, _ := range res.Headers {
	// 	valuePtrs[i] = &values[i]
	// }

	// for rows.Next() {
	// 	oneRow := make([]string, len(res.Headers))	
	// 	err = rows.Scan(valuePtrs...)
	// 	if err != nil {
	// 		return QueryResult{}, fmt.Errorf("failed scan row from result of query: '%s' error: %s", q, err)
	// 	}

	// 	for i, raw := range values {
    //         if raw == nil {
    //             oneRow[i] = "<null>"
    //         } else {
    //             oneRow[i] = string(raw)
    //         }
    //     }
	// 	res.Data = append(res.Data, oneRow)
	// }
	// return res, nil
}

// it is subtracted 3 hours (database timezone settings set to UTC, but time is MSK)
func fixTimezone(nt *sql.NullTime) {
	if nt.Valid {
		nt.Time = nt.Time.Local().Add(-3 * time.Hour)
	}
}