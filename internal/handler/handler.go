package handler

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"gitlab.com/wb-dynamics/wb-tg-query-executor/internal/postgres"
	"gitlab.com/wb-dynamics/wb-tg-query-executor/internal/telegram"
)

func DoWork (db *sql.DB) {
    tasks, err := postgres.GetTaskList(db)
    if err != nil {
        log.Printf("could not get tasks: %s", err)
    }

	for _, t := range tasks {		
		if isItTimeToWork(t.Id, t.ScheduleCron, t.LastExecutionTs) {
			log.Printf("Handle task id: %d. desc: %s", t.Id, t.ChatDescribe)
			handleTask(db, t)

			t.LastExecutionTs = sql.NullTime{Time: time.Now(), Valid: true}
			postgres.UpdateLastExecutionTs(db, t)
			log.Printf("id: %d last_execution_ts updated to: %s\n", t.Id, t.LastExecutionTs.Time.String())
		}
	}
}

func isItTimeToWork(id uint64, cronStr string, lastExecTs sql.NullTime) bool {
	schedule, err := cron.ParseStandard(cronStr)
	if err != nil {
		log.Printf("id: %d Could not parse cron string '%s'. error: %s", id, cronStr, err)
		return false
	}
	var lastExecTime time.Time
	if lastExecTs.Valid {
		lastExecTime = lastExecTs.Time
	}
	
	// log.Printf("id: %d last: %s\n", id, lastExecTime.String())
	nextTime := schedule.Next(lastExecTime)
	// log.Printf("id: %d next: %s\n", id, nextTime.String())
	// log.Printf("id: %d now : %s\n", id, time.Now().String())
	now := time.Now()
	diff := nextTime.Sub(now)
	log.Printf("id: %d. Next run: %s left: %s", id, nextTime, diff.String())
	return time.Now().After(nextTime)
}

func handleTask(db *sql.DB, task postgres.Task) {
	log.Printf("Executing query...\n")
	resultPretty, err := postgres.ExecQuery(db, task.SqlQuery)
	log.Printf("Executing query... OK\n")	
	if err != nil {
		log.Println(err)
		return 
	}
	msg := makeMessage(task.ChatDescribe, resultPretty)
	log.Printf("Sending message...\n")
	telegram.SendMessage(task.BotToken, task.ChatId, msg)
	log.Printf("Sending message... OK\n")
}

func makeMessage(desc, queryRes string) string {
	head := fmt.Sprintf("%s\n", desc)
	body := fmt.Sprintf("<pre>%s</pre>", queryRes)
	return head + body
}

