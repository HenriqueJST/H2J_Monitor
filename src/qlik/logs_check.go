package qlik

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bitbucket.org/h2j/h2j-qlik-monitor/src/config"
	"bitbucket.org/h2j/h2j-qlik-monitor/src/utils"
	"github.com/fsnotify/fsnotify"
	_ "github.com/go-sql-driver/mysql"
)

type LogsCheck struct{}

type notificationError struct {
	File      string
	Message   string
	TimeStart string
	TimeEnd   string
	Saas      string
}

func (l *LogsCheck) notifyError(filename string, message string, timestart string, timeend string) {
	var nome_cliente string
	var data_hora string
	now := time.Now()
	data_hora = now.Format("2006-01-02 03:04:05")
	nome_cliente = config.MonitorConfig.Cliente
	//------- Start SQL conf
	db, err := sql.Open("mysql", "monitoramento:ut39RDbFSfDPQxCR@tcp(mysql.h2j.com.br)/monitoramento")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//------- Finish SQL conf

	ne := notificationError{
		File:      filename,
		Message:   message,
		TimeStart: timestart,
		TimeEnd:   timeend,
		Saas:      config.MonitorConfig.Saas,
	}

	// Preparar o statement SQL cpu and memory
	stmt_cpu_memory, _ := db.Prepare("INSERT INTO qlik_logs (NOME_CLIENTE, DATA_HORA, FILE_NAME, MESSAGE, TIMESTART, TIMEEND, SAAS) VALUES (?, ?, ?, ?, ?, ?, ?)")
	defer stmt_cpu_memory.Close()

	// Executar o insert
	stmt_cpu_memory.Exec(nome_cliente, data_hora, ne.File, ne.Message, ne.TimeStart, ne.TimeEnd, ne.Saas)

}

func (l *LogsCheck) checkScriptFile(filename string) (err error) {
	time.Sleep(20 * time.Millisecond)
	file, _ := os.Open(filename)
	// if err != nil {
	// 	return

	// }
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	lines := []string{}
	index := 0
	for fileScanner.Scan() {
		index++
		line := fileScanner.Text()
		line = line[1:]
		lines = append(lines, line)

		if strings.Contains(fileScanner.Text(), "Execution finished") {
			startIndex := index - 10
			if startIndex < 0 {
				startIndex = 0
			}
			timeend := index - 1
			if timeend < 0 {
				timeend = 0
			}

			timestart := lines[1]

			l.notifyError(filename, strings.Join(lines[startIndex:index], "\n"), timestart, strings.Join(lines[timeend:index], "\n"))
			break
		}

	}

	return
}

func (l *LogsCheck) getFieldPos(fieldname string, headers []string) (pos int) {
	pos = -1
	for index, header := range headers {
		if header == fieldname {
			pos = index
			return
		}
	}
	return
}

func (l *LogsCheck) checkLogFile(filename string) (err error) {
	time.Sleep(20 * time.Millisecond)
	file, _ := os.Open(filename)
	// if err != nil {
	// 	return
	// }
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	lines, _ := reader.ReadAll()
	// if err != nil {
	// 	return
	// }
	headers := lines[0]
	lines = lines[1:]
	posLevel := l.getFieldPos("Severity", headers)
	if posLevel == -1 {
		posLevel = l.getFieldPos("Level", headers)
	}
	posMessage := l.getFieldPos("Message", headers)
	if posLevel >= 0 && posMessage >= 0 {
		for _, line := range lines {
			if utils.ContainsStr([]string{"ERROR", "WARN", "FATAL"}, strings.ToUpper(line[posLevel])) {
				posHostname := l.getFieldPos("Hostname", line)
				posException := l.getFieldPos("Exception", line)
				posUser := l.getFieldPos("UserId", line)
				posAppID := l.getFieldPos("Application", line)
				if posAppID == -1 {
					posAppID = l.getFieldPos("AppId", line)
				}
				posAppName := l.getFieldPos("AppName", line)
				posTaskName := l.getFieldPos("TaskName", line)
				posTaskID := l.getFieldPos("TaskId", line)
				message := fmt.Sprintf("%s: %s", line[posLevel], line[posMessage])
				if posException >= 0 {
					message = message + fmt.Sprintf(" / Exception: %s", line[posHostname])
				}
				if posHostname >= 0 {
					message = message + fmt.Sprintf(" / Hostname: %s", line[posHostname])
				}
				if posUser >= 0 {
					message = message + fmt.Sprintf(" / User: %s", line[posUser])
				}
				if posAppID >= 0 {
					message = message + fmt.Sprintf(" / App ID: %s", line[posAppID])
				}
				if posAppName >= 0 {
					message = message + fmt.Sprintf(" / App Name: %s", line[posAppName])
				}
				if posTaskID >= 0 {
					message = message + fmt.Sprintf(" / Task ID: %s", line[posTaskID])
				}
				if posTaskName >= 0 {
					message = message + fmt.Sprintf(" / Task Name: %s", line[posTaskName])
				}
				l.notifyError(filename, message, "", "")
			}
		}
	}
	return
}

func (l *LogsCheck) CheckFile(filename string, op fsnotify.Op) (err error) {
	if op != fsnotify.Remove {
		dir := filepath.Dir(filename)
		parent := filepath.Base(dir)
		if parent == "Script" {
			err = l.checkScriptFile(filename)
		} else {
			err = l.checkLogFile(filename)
		}
	}
	return
}
