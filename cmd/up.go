package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/unitedsymme/mittnite/internal/config"
	"github.com/unitedsymme/mittnite/pkg/files"
	"github.com/unitedsymme/mittnite/pkg/pidfile"
	"github.com/unitedsymme/mittnite/pkg/probe"
	"github.com/unitedsymme/mittnite/pkg/proc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	DefaultAPIAddress = "unix:///var/run/mittnite.sock"
)

var (
	probeListenPort  int
	pidFile          string
	apiEnabled       bool
	apiListenAddress string
	keepRunning      bool
)

func init() {
	log.StandardLogger().ExitFunc = func(i int) {
		defer func() {
			_ = recover() // prevent from printing trace
		}()
		panic(fmt.Sprintf("exit %d", i))
	}
	rootCmd.AddCommand(up)
	up.PersistentFlags().IntVarP(&probeListenPort, "probe-listen-port", "p", 9102, "set the port to listen for probe requests")
	up.PersistentFlags().StringVarP(&pidFile, "pidfile", "", "", "write mittnites process id to this file")
	up.PersistentFlags().BoolVarP(&apiEnabled, "api", "", false, "enables the api for remote or cli controlling")
	up.PersistentFlags().StringVarP(&apiListenAddress, "api-listen-address", "", DefaultAPIAddress, fmt.Sprintf("listen address for the api. Defaults to %q", DefaultAPIAddress))
	up.PersistentFlags().BoolVarP(&keepRunning, "keep-running", "k", false, "keep mittnite running even if no job is running anymore")
}

var up = &cobra.Command{
	Use:   "up",
	Short: "Render config files, start probes and processes",
	Long:  "This sub-command renders the configuration files, starts the probes and launches all configured processes",
	RunE: func(cmd *cobra.Command, args []string) error {
		ignitionConfig := &config.Ignition{
			Probes: nil,
			Files:  nil,
			Jobs:   nil,
		}

		pidFileHandle := pidfile.New(pidFile)

		if err := pidFileHandle.Acquire(); err != nil {
			return fmt.Errorf("failed to write pid file to %q: %w", pidFile, err)
		}

		defer func() {
			if err := pidFileHandle.Release(); err != nil {
				log.Errorf("error while cleaning up the pid file: %s", err)
			}
		}()

		if err := ignitionConfig.GenerateFromConfigDir(configDir); err != nil {
			return fmt.Errorf("failed while trying to generate ignition config from dir '%s': %w", configDir, err)
		}

		if err := files.RenderFiles(ignitionConfig.Files); err != nil {
			return fmt.Errorf("failed while rendering files from ignition config, err: %w", err)
		}

		signals := make(chan os.Signal)
		signal.Notify(signals,
			syscall.SIGTERM,
			syscall.SIGINT,
		)

		readinessSignals := make(chan os.Signal, 1)
		probeSignals := make(chan os.Signal, 1)
		procSignals := make(chan os.Signal, 1)

		go func() {
			for s := range signals {
				log.Infof("received event %s", s.String())
				readinessSignals <- s
				probeSignals <- s
				procSignals <- s
			}
		}()

		probeHandler, _ := probe.NewProbeHandler(ignitionConfig)

		go func() {
			log.Infof("probeServer listens on port %d", probeListenPort)

			if err := probe.RunProbeServer(probeHandler, probeSignals, probeListenPort); err != nil {
				log.Fatalf("probe server stopped with error: %s", err)
			} else {
				log.Info("probe server stopped without error")
			}
		}()

		go proc.ReapChildren()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var api *proc.Api
		if apiEnabled {
			api = proc.NewApi(apiListenAddress)

			defer api.Shutdown()
		}

		runner := proc.NewRunner(ctx, api, keepRunning, ignitionConfig)

		if err := runner.Init(); err != nil {
			return fmt.Errorf("runner failed to initialize: %w", err)
		}

		go func() {
			// start the API BEFORE waiting for readiness signals, so that the API is available
			// even if we're still waiting on some probes to become ready
			if err := runner.StartAPI(); err != nil {
				log.Fatalf("error while starting API: %v", err)
			}
		}()

		if err := probeHandler.Wait(readinessSignals); err != nil {
			return fmt.Errorf("probe handler failed while waiting for readiness signals: %w", err)
		}

		go func() {
			<-procSignals
			cancel()
		}()

		if err := runner.Boot(); err != nil {
			log.WithError(err).Fatal("runner error'ed during initialization")
		} else {
			log.Info("initialization complete")
		}

		if err := runner.Run(); err != nil {
			log.WithError(err).Fatal("service runner stopped with error")
		} else {
			log.Print("service runner stopped without error")
		}

		return nil
	},
}
