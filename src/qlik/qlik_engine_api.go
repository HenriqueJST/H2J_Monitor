package qlik

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"bitbucket.org/h2j/h2j-qlik-monitor/src/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/qlik-oss/enigma-go"
)

func Engine_Getlist(engineHost string, enginePort string, userName string, userDirectory string) {
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

	// Read client and root certificates.
	exePath, _ := os.Executable()
	dir := filepath.Dir(exePath)
	certFile := dir + "/certificados/client.pem"
	keyFile := dir + "/certificados/client_key.pem"
	caFile := dir + "/certificados/root.pem"

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		fmt.Println("Failed to load client certificate", err)
		panic(err)
	}

	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		fmt.Println("Failed to read root certificate", err)
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup TLS configuration.
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
	}

	ctx := context.Background()
	url := fmt.Sprintf("wss://%s:%s/app", engineHost, enginePort)

	// Notice how the user and directory is passed using the 'X-Qlik-User' header.
	headers := make(http.Header, 1)
	headers.Set("X-Qlik-User", fmt.Sprintf("UserDirectory=%s; UserId=%s", userDirectory, userName))

	global, err := enigma.Dialer{TLSClientConfig: tlsConfig}.Dial(ctx, url, headers)
	if err != nil {
		fmt.Println("Could not connect", err)
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a signal channel to capture termination signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	defer global.DisconnectFromServer()

	doc, _ := global.GetDocListRaw(ctx)

	type DocList struct {
		DOC_NAME         string `json:"qDocName"`
		FILE_SIZE        string `json:"qFileSize"`
		DOC_ID           string `json:"qDocId"`
		LAST_RELOAD_TIME string `json:"qLastReloadTime"`
	}

	var DocListFUll []DocList
	json.Unmarshal([]byte(doc), &DocListFUll)

	// Listando os qDocId
	for _, doc := range DocListFUll {
		// Preparar o statement SQL aplicacoes qlik
		stmt_cpu_memory, _ := db.Prepare("INSERT INTO qlik_doc_list (NOME_CLIENTE, DATA_HORA, DOC_NAME, FILE_SIZE, DOC_ID, LAST_RELOAD_TIME) VALUES (?, ?, ?, ?, ?, ?)")
		defer stmt_cpu_memory.Close()

		// Executar o insert
		stmt_cpu_memory.Exec(nome_cliente, data_hora, doc.DOC_NAME, doc.FILE_SIZE, doc.DOC_ID, doc.LAST_RELOAD_TIME)
	}

}
