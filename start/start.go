package start

import (
	"fmt"
	"io"
	"os"
	"strconv"
	SYS "syscall"
	"time"

	log "github.com/cihub/seelog"
	"github.com/grindlemire/WellsFarGO/rest"
	"github.com/grindlemire/WellsFarGO/unifier"
	"github.com/joho/godotenv"
	DEATH "github.com/vrecan/death"
)

// Run starts the webserver
func Run() {
	var goRoutines []io.Closer
	death := DEATH.NewDeath(SYS.SIGINT, SYS.SIGTERM)

	err := godotenv.Load()
	if err != nil {
		log.Critical("Error loading .env file")
		os.Exit(1)
	}

	port, err := GetEnvInt("SERVER_PORT")
	if err != nil {
		log.Critical("Unable to parse port in .env file: ", err)
		os.Exit(1)
	}

	restService := rest.NewRestService(port)
	goRoutines = append(goRoutines, restService)
	restService.Start()

	dbFile := os.Getenv("DB_FILE")
	csvFile := os.Getenv("CSV_FILE")
	formatType := os.Getenv("FORMAT_TYPE")

	unifier, err := unifier.NewUnifier(dbFile, csvFile, formatType)
	if err != nil {
		log.Critical("Error initializing the unifier")
		os.Exit(1)
	}

	err = unifier.AddNewData()
	if err != nil {
		log.Critical("Could not add new data: ", err)
		os.Exit(1)
	}

	qTime, _ := time.Parse("01/02/2006", "07/31/2016")
	results, _ := unifier.QueryDate(qTime)
	fmt.Printf("RESULTS: %#v\n", len(results))

	death.WaitForDeath(goRoutines...)

}

// GetEnvInt gets an environment variable and returns it as an int
func GetEnvInt(key string) (val int, err error) {
	strVal := os.Getenv(key)
	val, err = strconv.Atoi(strVal)
	return val, err
}
