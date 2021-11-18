package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"testing"
	"text/template"
	"time"

	"github.com/ory/dockertest/v3"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DSN struct {
	Host     string
	Username string
	Password string
	DBname   string
	Port     string
}

type httpClient struct {
	parent http.Client
}

func (client *httpClient) sendJsonReq(method, url string, reqBody []byte) (res *http.Response, resBody []byte, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-type", "application/json")

	resp, err := client.parent.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	resBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, resBody, nil
}

func waitForDBMSAndCreateConfig(pool *dockertest.Pool, resource *dockertest.Resource, dsn DSN) (confPath string, cleaner func()) {
	attempt := 0
	ok := false
	for attempt < 20 {
		attempt++
		connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			dsn.Host, dsn.Username, dsn.Password, dsn.DBname, dsn.Port)
		db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
		if err != nil {
			log.Infof("[waitForDBMSAndCreateConfig] gorm.Open failed: %v, waiting... (attempt %d)", err, attempt)
			time.Sleep(1 * time.Second)
			continue
		}

		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
		ok = true
		break
	}

	if !ok {
		_ = pool.Purge(resource)
		log.Panicf("[waitForDBMSAndCreateConfig] couldn't connect to PostgreSQL")
	}

	tmpl, err := template.New("config").Parse(`
server_mode: debug
bind_ip: 0.0.0.0
http_port: 8080
postgres:
  host: {{.Host}}
  username: {{.Username}}
  password: {{.Password}}
  dbname: {{.DBname}}
  port: {{.Port}}
`)
	if err != nil {
		_ = pool.Purge(resource)
		log.Panicf("[waitForDBMSAndCreateConfig] template.Parse failed: %v", err)
	}

	var configBuff bytes.Buffer
	err = tmpl.Execute(&configBuff, dsn)
	if err != nil {
		_ = pool.Purge(resource)
		log.Panicf("[waitForDBMSAndCreateConfig] tmpl.Execute failed: %v", err)
	}

	confFile, err := ioutil.TempFile("", "app.*.yaml")
	if err != nil {
		_ = pool.Purge(resource)
		log.Panicf("[waitForDBMSAndCreateConfig] ioutil.TempFile failed: %v", err)
	}
	log.Infof("[waitForDBMSAndCreateConfig] confFile.Name = %s", confFile.Name())

	_, err = confFile.WriteString(configBuff.String())
	if err != nil {
		_ = pool.Purge(resource)
		log.Panicf("[waitForDBMSAndCreateConfig] confFile.WriteString failed: %v", err)
	}

	err = confFile.Close()
	if err != nil {
		_ = pool.Purge(resource)
		log.Panicf("[waitForDBMSAndCreateConfig] confFile.Close failed: %v", err)
	}

	cleanerFunc := func() {
		err := pool.Purge(resource)
		if err != nil {
			log.Panicf("[waitForDBMSAndCreateConfig] pool.Purge failed: %v", err)
		}

		err = os.Remove(confFile.Name())
		if err != nil {
			log.Panicf("[waitForDBMSAndCreateConfig] os.Remove failed: %v", err)
		}
	}
	return confFile.Name(), cleanerFunc
}

func startPostgreSQL() (confPath string, cleaner func()) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Panicf("[startPostgreSQL] dockertest.NewPool failed: %v", err)
	}

	resource, err := pool.Run(
		"postgres", "11",
		[]string{
			"POSTGRES_DB=restservice",
			"POSTGRES_PASSWORD=this_is_postgres",
		},
	)
	if err != nil {
		log.Panicf("[startPostgreSQL] pool.Run failed: %v", err)
	}

	host, port, _ := net.SplitHostPort(resource.GetHostPort("5432/tcp"))
	dsn := DSN{
		Host:     host,
		Username: "postgres",
		Password: "this_is_postgres",
		DBname:   "postgres",
		Port:     port,
	}
	return waitForDBMSAndCreateConfig(pool, resource, dsn)
}

func TestMain(m *testing.M) {
	log.Infoln("[TestMain] About to start PostgreSQL...")
	confPath, stopDB := startPostgreSQL()
	log.Infoln("[TestMain] PostgreSQL started!")

	cmd := exec.Command("./simple-rest", "-c", confPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		stopDB()
		log.Panicf("[TestMain] cmd.Start failed: %v", err)
	}
	log.Infof("[TestMain] cmd.Process.Pid = %d", cmd.Process.Pid)

	attempt := 0
	ok := false
	client := httpClient{}
	for attempt < 20 {
		attempt++
		_, _, err := client.sendJsonReq("GET", "http://localhost:8080/records/0", []byte{})
		if err != nil {
			log.Infof("[TestMain] client.sendJsonReq failed: %v, waiting... (attempt %d)", err, attempt)
			time.Sleep(1 * time.Second)
			continue
		}

		ok = true
		break
	}

	if !ok {
		stopDB()
		_ = cmd.Process.Kill()
		log.Panicf("[TestMain] REST API is unavailable")
	}

	log.Infoln("[TestMain] REST API ready! Executing m.Run()")
	code := m.Run()

	log.Infoln("[TestMain] Cleaning up...")
	_ = cmd.Process.Signal(syscall.SIGTERM)
	stopDB()
	os.Exit(code)
}

func TestCRUD(t *testing.T) {
	t.Parallel()

	type PhoneRecord struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}
	client := httpClient{}

	respID := struct {
		ID int
	}{}

	// CREATE
	record := PhoneRecord{
		Name:  "Alice",
		Phone: "123",
	}
	httpBody, err := json.Marshal(record)
	require.NoError(t, err)
	resp, respBody, err := client.sendJsonReq("POST", "http://localhost:8080/records", httpBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	err = json.Unmarshal(respBody, &respID)
	require.NoError(t, err)
	require.NotEqual(t, 0, respID.ID)

	// READ
	respRecord := PhoneRecord{}
	resp, respBody, err = client.sendJsonReq("GET", "http://localhost:8080/records/"+strconv.Itoa(respID.ID), []byte{})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	err = json.Unmarshal(respBody, &respRecord)
	require.NoError(t, err)
	require.Equal(t, respRecord.ID, respID.ID)
	require.Equal(t, respRecord.Name, record.Name)
	require.Equal(t, respRecord.Phone, record.Phone)

	// UPDATE
	record.Name = "John"
	record.ID = respID.ID
	httpBody, err = json.Marshal(record)
	require.NoError(t, err)
	resp, _, err = client.sendJsonReq("PUT", "http://localhost:8080/records/"+strconv.Itoa(record.ID), httpBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp, respBody, err = client.sendJsonReq("GET", "http://localhost:8080/records/"+strconv.Itoa(record.ID), []byte{})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	err = json.Unmarshal(respBody, &respRecord)
	require.NoError(t, err)
	require.Equal(t, respRecord.Name, record.Name)

	// DELETE

}
