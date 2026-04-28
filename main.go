package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/tinydarkforge/intake/app"
	"github.com/tinydarkforge/intake/services"
	"github.com/tinydarkforge/intake/types"
	"github.com/tinydarkforge/intake/ui"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	repoRe = regexp.MustCompile(`^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+$`)
)

func main() {
	var (
		flagVersion = flag.Bool("version", false, "print version")
		flagRepo    = flag.String("repo", "", "GitHub repo (owner/name)")
		flagModel   = flag.String("model", "", "Ollama model")
		flagHost    = flag.String("ollama-host", "", "Ollama host URL")
		flagNoSound = flag.Bool("no-sound", false, "Disable sound")
		flagDebug   = flag.Bool("debug", false, "Debug mode")
	)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "intake — agentic GitHub issue creator")
		fmt.Fprintln(os.Stderr, "Usage: intake [flags]")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nEnv vars: INTAKE_REPO, INTAKE_MODEL, OLLAMA_HOST, INTAKE_TIMEOUT, INTAKE_MAX_TURNS, INTAKE_TEMPLATE_DIR, INTAKE_DEBUG")
	}
	flag.Parse()

	if *flagVersion {
		fmt.Printf("intake %s (%s) %s\n", version, commit, date)
		os.Exit(0)
	}

	cfg, err := types.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: config load: %v\n", err)
	}

	// CLI flags override config + env.
	if *flagRepo != "" {
		cfg.Repo = *flagRepo
	}
	if *flagModel != "" {
		cfg.Model = *flagModel
	}
	if *flagHost != "" {
		cfg.OllamaHost = *flagHost
	}
	if *flagNoSound {
		cfg.Sound = false
	}
	if *flagDebug {
		cfg.Debug = true
	}

	if cfg.Repo != "" && !repoRe.MatchString(cfg.Repo) {
		fmt.Fprintf(os.Stderr, "intake: invalid repo %q — expected owner/name\n", cfg.Repo)
		os.Exit(1)
	}

	// Build services.
	gh := services.NewGitHub(cfg.Repo)
	ollama := services.NewOllama(cfg.OllamaHost, cfg.Model, time.Duration(cfg.TimeoutSec)*time.Second)

	{
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := ollama.Ping(ctx); err != nil {
			fmt.Fprintln(os.Stderr, "warning: ollama unreachable — create flow will fail until ollama is running")
		}
	}

	tmpls, err := services.LoadTemplates(cfg.TemplateDir)
	if err != nil || len(tmpls) == 0 {
		tmpls = []types.Template{{Name: "Default", Body: "## Description\n\n<!-- fill in -->"}}
	}

	// Build screen models.
	localModel := ui.NewLocalPane()
	listModel := ui.NewList(gh)
	createModel := ui.NewCreate(gh, ollama, tmpls, cfg)
	detailModel := ui.NewDetail(gh)
	settingsModel := ui.NewSettings(cfg)

	root := app.New()
	root.Cfg = cfg
	root.Local = &localModel
	root.List = &listModel
	root.Create = &createModel
	root.Detail = &detailModel
	root.Settings = &settingsModel

	p := tea.NewProgram(
		root,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Give create model a reference to the program for Send-based streaming.
	createModel = createModel.SetProgram(p)
	root.Create = &createModel

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "intake: %v\n", err)
		os.Exit(1)
	}
}
