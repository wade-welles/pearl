package cmd

import (
	"os"

	"github.com/inconshreveable/log15"
	"github.com/mmcloughlin/pearl"
	"github.com/mmcloughlin/pearl/log"
	"github.com/mmcloughlin/pearl/meta"
	"github.com/mmcloughlin/pearl/torconfig"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a relay server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return serve()
	},
}
var (
	nickname string
	port     int
	logfile  string
)

func init() {
	serveCmd.Flags().StringVarP(&nickname, "nickname", "n", "pearl", "nickname")
	serveCmd.Flags().IntVarP(&port, "port", "p", 9111, "relay port")
	serveCmd.Flags().StringVarP(&logfile, "logfile", "l", "pearl.json", "log file")

	rootCmd.AddCommand(serveCmd)
}

func serve() error {
	config := &torconfig.Config{
		Nickname: nickname,
		ORPort:   uint16(port),
		Platform: meta.Platform.String(),
		Contact:  "https://github.com/mmcloughlin/pearl",
	}

	base := log15.New()
	fh, err := log15.FileHandler(logfile, log15.JsonFormat())
	if err != nil {
		return err
	}
	base.SetHandler(log15.MultiHandler(
		log15.LvlFilterHandler(log15.LvlInfo,
			log15.StreamHandler(os.Stdout, log15.TerminalFormat()),
		),
		fh,
	))
	logger := log.NewLog15(base)

	r, err := pearl.NewRouter(config, logger)
	if err != nil {
		return err
	}

	go func() {
		r.Serve()
	}()

	authority := "127.0.0.1:7000"
	desc, err := r.Descriptor()
	if err != nil {
		return err
	}
	err = desc.PublishToAuthority(authority)
	if err != nil {
		return err
	}
	logger.With("authority", authority).Info("published descriptor")

	select {}
}
