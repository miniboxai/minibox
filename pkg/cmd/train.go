package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	bopt "minibox.ai/pkg/backend/option"

	"minibox.ai/pkg/backend/docker"
	"minibox.ai/pkg/core"

	pcontext "minibox.ai/pkg/core/context"
	"minibox.ai/pkg/core/job"
	"minibox.ai/pkg/core/option"
	"minibox.ai/pkg/logger"
)

var (
	configFile string
	backend    string
	wait       bool
)

var trainCmd = &cobra.Command{
	Use:   "train",
	Short: "train start, create a training job",
	Run: func(cmd *cobra.Command, args []string) {
		if debug {
			// setLoggerLevel(zapcore.ErrorLevel)
			enableDebug()
		}

		if err := LoadAuth(startTrainExec); err != nil {
			fmt.Printf("auth error you must do `mini login` in terminal")
		}
	},
}

func startTrainExec(token string) error {
	logger.S().Info("Train starting...")
	f := openConfig()
	defer f.Close()

	var rootPath, prjName string
	dir, _ := os.Getwd()
	rootPath = dir

	prjName = filepath.Base(rootPath)

	prj, err := core.FromConfig(f, &option.ConfigOpt{
		ProjectName: prjName,
		RootPath:    rootPath,
	})

	if err != nil {
		logger.S().Errorf("parse config error %s", err)
		return nil
	}

	var initialParams = map[string]interface{}{
		"projectName": prj.Name,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	ctx = context.WithValue(
		ctx,
		pcontext.ContextKey,
		pcontext.NewContext(initialParams))

	job, err := prj.Start(ctx, &option.StartProjectOption{
		ProjectName: prj.Name,
		Wait:        wait,
	})

	if err != nil {
		logger.S().Errorf("project Start error: %s", err)
		return nil
	}

	logger.S().Infof("job: %#v", job)

	executeJob(job, ctx, &bopt.ExecuteOption{
		WaitStatus: wait,
		Log:        wait,
	})

	return nil
}

func openConfig() *os.File {
	f, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("open project config failed, %s", err)
	}
	return f
	// defer f.Close()
}

func cwd() string {
	if curDir, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		return curDir
	}
}

func executeJob(j *job.Job, ctx context.Context, opt *bopt.ExecuteOption) error {

	if backend == "docker" {
		if b, err := docker.NewBackend(); err != nil {
			panic(err)
		} else {
			b.Execute(j, ctx, opt)

			if opt.Log {
				go io.Copy(os.Stdout, j.Output())
			}

			if opt.WaitStatus {
				j.Wait()
			}

		}
	} else {
		panic(fmt.Errorf("we dont no yeah impo this Backend %s", backend))
	}
	return nil
}

func init() {
	flag := trainCmd.PersistentFlags()
	flag.StringVar(&configFile, "config", "minibox.yaml", "Project's configure file.")
	flag.BoolVarP(&debug, "debug", "d", false, "Print all debug information")
	flag.BoolVarP(&wait, "wait", "w", false, "Wait for Container exit")
	flag.StringVarP(&backend, "backend", "b", "docker", "Use [docker, kuberneters, minibox] different Backend system ")
}
