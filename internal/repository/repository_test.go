package repository

import (
	"database/sql"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/ory/dockertest/v3"
	"os"
	"testing"
)

//func TestRespondsWithLove(t *testing.T) {
//
//	pool, err := dockertest.NewPool("")
//	require.NoError(t, err, "could not connect to Docker")
//
//	resource, err := pool.Run("docker-gs-ping", "latest", []string{})
//	require.NoError(t, err, "could not start container")
//
//	t.Cleanup(func() {
//		require.NoError(t, pool.Purge(resource), "failed to remove container")
//	})
//
//	var resp *http.Response
//
//	err = pool.Retry(func() error {
//		resp, err = http.Get(fmt.Sprint("http://localhost:", resource.GetPort("8080/tcp"), "/"))
//		if err != nil {
//			t.Log("container not ready, waiting...")
//			return err
//		}
//		return nil
//	})
//	require.NoError(t, err, "HTTP error")
//	defer resp.Body.Close()
//
//	require.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")
//
//	body, err := io.ReadAll(resp.Body)
//	require.NoError(t, err, "failed to read HTTP body")
//
//	// Finally, test the business requirement!
//	require.Contains(t, string(body), "<3", "does not respond with love?")
//}



var db *sql.DB

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("postgres", "latest", []string{"POSTGRES_PASSWORD=secret"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestSomething(t *testing.T) {
	// db.Query()
}


//func StartPostgreSQL() (confPath string, cleaner func()) {
//	pool, err := dockertest.NewPool("")
//	if err != nil {
//		log.Panicf("dockertest.NewPool failed: %v", err)
//	}
//
//	resource, err := pool.Run(
//		"postgres", "11",
//		[]string{
//			"POSTGRES_DB=restservice",
//			"POSTGRES_PASSWORD=s3cr3t",
//		},
//	)
//	if err != nil {
//		log.Panicf("pool.Run failed: %v", err)
//	}
//
//	connString := "postgres://postgres:s3cr3t@"+
//		resource.GetHostPort("5432/tcp")+
//		"/restservice?sslmode=disable"
//	attempt := 0
//	ok := false
//	for attempt < 20 {
//		attempt++
//		conn, err := pgx.Connect(context.Background(), connString)
//		if err != nil {
//			log.Infof("pgx.Connect failed: %v, waiting... (attempt %d)",
//				err, attempt)
//			time.Sleep(1 * time.Second)
//			continue
//		}
//
//		_ = conn.Close(context.Background())
//		ok = true
//		break
//	}
//
//	if !ok {
//		_ = pool.Purge(resource)
//		log.Panicf("Couldn't connect to PostgreSQL")
//	}
//
//	tmpl, err := template.New("config").Parse(`
//loglevel: debug
//listen: 0.0.0.0:8080
//db:
//  url: {{.ConnString}}
//`)
//	if err != nil {
//		_ = pool.Purge(resource)
//		log.Panicf("template.Parse failed: %v", err)
//	}
//
//	configArgs := struct {
//		ConnString string
//	} {
//		ConnString: connString,
//	}
//	var configBuff bytes.Buffer
//	err = tmpl.Execute(&configBuff, configArgs)
//	if err != nil {
//		_ = pool.Purge(resource)
//		log.Panicf("tmpl.Execute failed: %v", err)
//	}
//
//	confFile, err := ioutil.TempFile("", "config.*.yaml")
//	if err != nil {
//		_ = pool.Purge(resource)
//		log.Panicf("ioutil.TempFile failed: %v", err)
//	}
//
//	log.Infof("confFile.Name = %s", confFile.Name())
//
//	_, err = confFile.WriteString(configBuff.String())
//	if err != nil {
//		_ = pool.Purge(resource)
//		log.Panicf("confFile.WriteString failed: %v", err)
//	}
//
//	err = confFile.Close()
//	if err != nil {
//		_ = pool.Purge(resource)
//		log.Panicf("confFile.Close failed: %v", err)
//	}
//
//	cleanerFunc := func() {
//		// purge the container
//		err := pool.Purge(resource)
//		if err != nil {
//			log.Panicf("pool.Purge failed: %v", err)
//		}
//
//		err = os.Remove(confFile.Name())
//		if err != nil {
//			log.Panicf("os.Remove failed: %v", err)
//		}
//	}
//
//	return confFile.Name(), cleanerFunc
//}
//
//func TestMain(m *testing.M) {
//	fmt.Println("About to start PostgreSQL...")
//	confPath, stopPostgreSQL := StartPostgreSQL()
//	fmt.Println("PostgreSQL started!")
//
//	fmt.Println(confPath)
//	stopPostgreSQL()
//}