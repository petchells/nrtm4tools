package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/service"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")
var notificationURL = flag.String("url", "", "URL to notification JSON")
var sourceName = flag.String("source", "", "The name of the source")
var sourceLabel = flag.String("label", "", "The label for the source. Can be empty.")

func main() {
	envVars := []string{"PG_DATABASE_URL", "NRTM4_FILE_PATH"}
	for _, ev := range envVars {
		if len(os.Getenv(ev)) <= 0 {
			log.Fatalln("Environment variable not set: ", ev)
		}
	}
	dbURL := os.Getenv("PG_DATABASE_URL")
	boltDBPath := os.Getenv("BOLT_DATABASE_PATH")
	nrtmFilePath := os.Getenv("NRTM4_FILE_PATH")
	config := service.AppConfig{
		NRTMFilePath:     nrtmFilePath,
		PgDatabaseURL:    dbURL,
		BoltDatabasePath: boltDBPath,
	}
	commander := nrtm4.InitializeCommandProcessor(config)
	connectCommand := func(args []string) {
		// A real program (not an example) would use flag.ExitOnError.
		fs := flag.NewFlagSet("connect", flag.ExitOnError)
		if err := fs.Parse(args); err != nil {
			fmt.Printf("error: %s", err)
			return
		}
		if len(*notificationURL) == 0 {
			log.Fatal("URL must be provided")
		}
		commander.Connect(*notificationURL, *sourceLabel)
	}

	updateCommand := func(args []string) {
		fs := flag.NewFlagSet("update", flag.ExitOnError)
		src := fs.String("source", "", "The name of the source")
		lbl := fs.String("label", "", "The label for the source. Can be empty.")
		if err := fs.Parse(args); err != nil {
			fmt.Printf("error: %s", err)
			return
		}
		if len(*src) == 0 {
			log.Fatalf("Source name must be provided")
		}
		commander.Update(*src, *lbl)
	}

	listCommand := func(args []string) {
		fs := flag.NewFlagSet("list", flag.ExitOnError)
		src := fs.String("source", "", "The name of the source")
		lbl := fs.String("label", "", "The label for the source. Can be empty.")
		if err := fs.Parse(args); err != nil {
			fmt.Printf("error: %s", err)
			return
		}
		commander.ListSources(*src, *lbl)
	}

	runCmd := func(args []string) {
		if len(args) >= 2 {
			subArgs := args[2:]
			switch args[1] {
			case "connect":
				connectCommand(subArgs)
				return
			case "update":
				updateCommand(subArgs)
				return
			case "list":
				listCommand(subArgs)
				return
			default:
			}
		}
		log.Print(usage())
		flag.Usage()
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}

	runCmd(os.Args)
}

func usage() string {
	return `
	%v <command> OPTIONS

	command: [connect|update]

	The client reads two properties from environment variables, which must be set:

	PG_DATABASE_URL
	URL to the PostgreSQL database in this format:
	postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]

	NRTM4_FILE_PATH
	The path where downloaded NRTMv4 snapshot and delta files will be written. If
	the directory does not exist it will be created. The files are only needed
	during updates; when the update is complete the files can be removed.
	...Which is probably a good idea, there's a lot of files.


	E.g.
	envvars="\
	PG_DATABASE_URL=postgres://nrtm4:nrtm4@localhost:5432/nrtm4?sslmode=disable \
	NRTM4_FILE_PATH=/tmp/nrtm4"

	env ${envvars} nrtm4client connect -url https://nrtm4.example.zz/notification.json

	env ${envvars} nrtm4client update -source EXAMPLE
	`
}
