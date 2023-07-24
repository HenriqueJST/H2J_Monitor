package server

import (
	"database/sql"
	"time"

	"bitbucket.org/h2j/h2j-qlik-monitor/src/config"
	"bitbucket.org/h2j/h2j-qlik-monitor/src/host"
	_ "github.com/go-sql-driver/mysql"
)

func Server_Service() {
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

	p := host.GetInfoProcess()
	d := host.GetDisks()
	h := host.GetHostInfo()

	for _, doc := range p {
		// Preparar o statement SQL process
		stmt_process, _ := db.Prepare("INSERT INTO server_process (NOME_CLIENTE, DATA_HORA, NOME_PROCESSO, USO_CPU, USO_MEMORIA, ISRUNNING) VALUES (?, ?, ?, ?, ?, ?)")
		defer stmt_process.Close()

		// Executar o insert
		stmt_process.Exec(nome_cliente, data_hora, doc.Name, doc.CPUPercent, doc.MemoryPercent, doc.IsRunning)

	}
	for _, doc := range d {
		// Preparar o statement SQL disks
		stmt_disks, _ := db.Prepare("INSERT INTO server_disks (NOME_CLIENTE, DATA_HORA, DEVICE, FSTYPE, TOTALSPACE, USEDSPACE, FREESPACE) VALUES (?, ?, ?, ?, ?, ?, ?)")
		defer stmt_disks.Close()
		// Executar o insert
		stmt_disks.Exec(nome_cliente, data_hora, doc.Device, doc.FSType, doc.TotalSpace, doc.UsedSpace, doc.FreeSpace)
	}

	for _, doc := range h {
		// Preparar o statement SQL cpu and memory
		stmt_cpu_memory, _ := db.Prepare("INSERT INTO server_info (NOME_CLIENTE, DATA_HORA, CPU, MEMORY, HOSTNAME) VALUES (?, ?, ?, ?, ?)")
		defer stmt_cpu_memory.Close()

		// Executar o insert
		stmt_cpu_memory.Exec(nome_cliente, data_hora, doc.CPU, doc.Memory, doc.HostName)
	}
}
