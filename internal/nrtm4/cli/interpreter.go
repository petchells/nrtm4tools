package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/service"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile = flag.String("memprofile", "", "write memory profile to this file")
)

const mandatorySourceMessage = "Source name must be provided with the -source flag"

// Exec reads the command line args and invokes functions on the commander
func Exec(commander service.CommandExecutor) {

	connectCommand := func(args []string) {
		fs := flag.NewFlagSet("connect", flag.ExitOnError)
		notificationURL := fs.String("url", "", "URL to notification JSON")
		sourceLabel := fs.String("label", "", "The label for the source. Can be empty.")
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
			log.Fatalf(mandatorySourceMessage)
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

	replaceLabelCommand := func(args []string) {
		fs := flag.NewFlagSet("rename", flag.ExitOnError)
		src := fs.String("source", "", "The name of the source")
		lbl := fs.String("label", "", "The label for the source. Can be empty.")
		tolbl := fs.String("to", "", "The replacement label text")
		if err := fs.Parse(args); err != nil {
			fmt.Printf("error: %s", err)
			return
		}
		if len(*src) == 0 {
			log.Fatalf(mandatorySourceMessage)
		}
		if len(*lbl) == 0 && len(*tolbl) == 0 {
			log.Fatalf("At least -label or -to must be specified")
		}
		commander.ReplaceLabel(*src, *lbl, *tolbl)
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
			case "rename":
				replaceLabelCommand(subArgs)
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
