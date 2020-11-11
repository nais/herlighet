package main

import (
    "errors"
    "strconv"
    "os"
    "database/sql"
    _ "go.uber.org/zap"
    "fmt"
    _ "github.com/lib/pq"
)

const (
  host     = "localhost"
  port     = 5433
  user     = "herlighet"
  password = "herlighet"
  confdb   = "herlighet"
)

func openDb() *sql.DB {
    psqlInfo := getDbConfString()
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        panic(err)
    }
    err = db.Ping()
    if err != nil {
        panic(err)
    }
    return db
}

func smellTest() {
    var count int
    db := openDb()
    row := db.QueryRow("SELECT COUNT(*) FROM databases;")
    err := row.Scan(&count)
    if err != nil {
        panic(err)
    }
    logger.Info(fmt.Sprintf("%d databases in access list", count))
    db.Close()
}

func getDbConfField(key string, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    os.Setenv(key, fallback)
    logger.Info(fmt.Sprintf(`"%s" is not defined; using fallback value.`, key))
    return fallback
}

func getDbConfString() string {
    user := getDbConfField("HERLIGHET_DBUSER", "herlighet")
    pass := getDbConfField("HERLIGHET_DBPASS", "herlighet")
    host := getDbConfField("HERLIGHET_DBHOST", "localhost")
    port, _ := strconv.Atoi(getDbConfField("HERLIGHET_DBPORT", "5433"))
    dbname := getDbConfField("HERLIGHET_DBNAME", "herlighet")
    logger.Info(fmt.Sprintf("Connecting as %s to %s:%d/%s", user, host, port, dbname))
    return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, pass, dbname)
}

func getRearEnd(dbname string) (string, error) {
    db := openDb()
    defer db.Close()

    sqlStatement := `SELECT hostname, naisdevice FROM databases WHERE dbname=$1;`

    var hostname string
    var naisdevice bool

    row := db.QueryRow(sqlStatement, dbname)
    switch err := row.Scan(&hostname, &naisdevice); err {
    case sql.ErrNoRows:
        return "", errors.New("naisdevice access not configured for this database; please ask #postgres-på-laptop for help with updating database-iac.")
    case nil:
        if !naisdevice {
            return "", errors.New("naisdevice access disabled for this database; please update database-iac.")
        } else {
            return hostname, nil
        }
    default:
        return "", err
    }
}
