package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strconv"
	"strings"

	"docflow/internal/buildinfo"
	"docflow/internal/db"
	"docflow/internal/engine"
	"docflow/internal/logger"
	"docflow/internal/model"
	"docflow/internal/validator"
	"docflow/pkg/config"
	"docflow/pkg/governance"
	"docflow/pkg/index"
	"docflow/pkg/pattern"
	"docflow/pkg/rules"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const defaultDB = "docflow.db"

var (
	dbPath       string
	configFile   string
	logLevel     string
	logFormat    string
	cpuProfile   string
	memProfile   string
	cacheEnabled = true
	noCache      bool
	cacheDir     string
	changedOnly  bool
	sinceRef     string
)

func Execute() error {
	rootCmd := &cobra.Command{
		Use:           "docflow",
		Short:         "docflow – CLI do zarządzania dokumentacją i szablonami",
		Long:          "Framework do zarządzania dokumentacją z metadanymi.\nSkanuje repozytorium, waliduje metadane, buduje graf zależności.",
		Version:       buildinfo.Version,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cpuProfile != "" {
				f, err := os.Create(cpuProfile)
				if err != nil {
					return err
				}
				if err := pprof.StartCPUProfile(f); err != nil {
					return err
				}
			}
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if cpuProfile != "" {
				pprof.StopCPUProfile()
			}
			if memProfile != "" {
				f, err := os.Create(memProfile)
				if err == nil {
					_ = pprof.WriteHeapProfile(f)
					f.Close()
				}
			}
		},
	}
	rootCmd.SetVersionTemplate("{{.Name}} " + buildinfo.FullVersion() + "\n")

	rootCmd.PersistentFlags().StringVar(&dbPath, "db", defaultDB, "ścieżka do pliku bazy SQLite")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "ścieżka do pliku konfiguracji (domyślnie: docflow.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "poziom logowania (debug, info, warn, error)")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "", "format logów (text|json)")
	rootCmd.PersistentFlags().StringVar(&cpuProfile, "cpu-profile", "", "zapisz profil CPU do pliku")
	rootCmd.PersistentFlags().StringVar(&memProfile, "mem-profile", "", "zapisz profil pamięci do pliku")
	rootCmd.PersistentFlags().BoolVar(&cacheEnabled, "cache", true, "włącz cache execution layer (SQLite)")
	rootCmd.PersistentFlags().BoolVar(&noCache, "no-cache", false, "wyłącz cache execution layer")
	rootCmd.PersistentFlags().StringVar(&cacheDir, "cache-dir", "", "katalog cache execution layer (domyślnie z config: .docflow/cache)")
	rootCmd.PersistentFlags().BoolVar(&changedOnly, "changed-only", false, "parsuj tylko zmienione pliki, ale ewaluuj pełny zbiór facts")
	rootCmd.PersistentFlags().StringVar(&sinceRef, "since", "", "git ref do wyznaczania zmienionych plików (np. origin/main)")

	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(importCmd())
	rootCmd.AddCommand(phaseCmd())
	rootCmd.AddCommand(taskCmd())
	rootCmd.AddCommand(statusCmd())
	rootCmd.AddCommand(scanCmd())
	rootCmd.AddCommand(validateCmd())
	rootCmd.AddCommand(planCmd())
	rootCmd.AddCommand(analyzePatternsCmd())
	rootCmd.AddCommand(statsCmd())
	rootCmd.AddCommand(validateTemplatesCmd())
	rootCmd.AddCommand(findDuplicatesCmd())
	rootCmd.AddCommand(recommendCmd())
	rootCmd.AddCommand(setsCmd())
	rootCmd.AddCommand(templatesCmd())
	rootCmd.AddCommand(templateImpactCmd())
	rootCmd.AddCommand(graphCmd())
	rootCmd.AddCommand(baselineCmd())
	rootCmd.AddCommand(changesCmd())
	rootCmd.AddCommand(migrateSectionsCmd())
	rootCmd.AddCommand(complianceCmd())
	rootCmd.AddCommand(paletteCmd())
	rootCmd.AddCommand(healthCmd())
	rootCmd.AddCommand(fixCmd())

	return rootCmd.Execute()
}

func openStore() (*db.Store, error) {
	return db.New(dbPath)
}

func initCmd() *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Inicjalizuje nowy projekt",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := openStore()
			if err != nil {
				return err
			}
			defer store.Close()

			if err := store.InitProject(name); err != nil {
				return err
			}

			fmt.Printf("Projekt '%s' zainicjalizowany (baza: %s)\n", name, dbPath)
			return nil
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "docflow", "nazwa projektu")
	return cmd
}

func importCmd() *cobra.Command {
	var filePath string
	var clear bool
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Importuje plan z pliku YAML lub JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := openStore()
			if err != nil {
				return err
			}
			defer store.Close()

			data, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("cannot read file %s: %w", filePath, err)
			}

			var plan model.PlanImport

			if strings.HasSuffix(filePath, ".json") {
				err = json.Unmarshal(data, &plan)
			} else {
				err = yaml.Unmarshal(data, &plan)
			}
			if err != nil {
				return fmt.Errorf("cannot parse file: %w", err)
			}

			if plan.ProjectName != "" {
				if err := store.InitProject(plan.ProjectName); err != nil {
					return err
				}
			}

			if clear {
				if err := store.ClearPhases(); err != nil {
					return err
				}
			}

			imported := 0
			for _, p := range plan.Phases {
				if p.Status == "" {
					p.Status = model.StatusTodo
				}
				_, err := store.AddPhase(p)
				if err != nil {
					fmt.Fprintf(os.Stderr, "WARN: nie udało się dodać fazy day=%d: %v\n", p.Day, err)
					continue
				}
				imported++
			}

			for _, p := range plan.Phases {
				if len(p.DependsOn) == 0 {
					continue
				}
				missing, err := store.LinkDependencies(p.Day, p.DependsOn)
				if err != nil {
					fmt.Fprintf(os.Stderr, "WARN: zależności day=%d: %v\n", p.Day, err)
				}
				for _, m := range missing {
					fmt.Fprintf(os.Stderr, "WARN: faza day=%d zależy od day=%d, która nie istnieje\n", p.Day, m)
				}
			}

			fmt.Printf("Zaimportowano %d faz z %s\n", imported, filePath)
			return nil
		},
	}
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "plik z planem (YAML lub JSON)")
	cmd.Flags().BoolVar(&clear, "clear", false, "wyczyść istniejące fazy przed importem")
	cmd.MarkFlagRequired("file")
	return cmd
}

func phaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "phase",
		Short: "Zarządzanie fazami projektu",
	}
	cmd.AddCommand(phaseListCmd())
	cmd.AddCommand(phaseShowCmd())
	cmd.AddCommand(phaseUpdateCmd())
	cmd.AddCommand(phaseAddCmd())
	return cmd
}

func phaseListCmd() *cobra.Command {
	var statusFilter string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lista faz",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := openStore()
			if err != nil {
				return err
			}
			defer store.Close()

			phases, err := store.ListPhases(statusFilter)
			if err != nil {
				return err
			}

			if len(phases) == 0 {
				fmt.Println("Brak faz. Użyj 'docflow import -f plan.yaml' aby zaimportować plan.")
				return nil
			}

			fmt.Printf("%-6s %-12s %-40s %s\n", "DAY", "STATUS", "NAZWA", "ZALEŻNOŚCI")
			fmt.Println(strings.Repeat("─", 80))
			for _, p := range phases {
				statusIcon := statusToIcon(p.Status)
				deps := ""
				if len(p.DependsOn) > 0 {
					depStrs := make([]string, len(p.DependsOn))
					for i, d := range p.DependsOn {
						depStrs[i] = fmt.Sprintf("d%d", d)
					}
					deps = strings.Join(depStrs, ", ")
				}
				name := p.Name
				if len(name) > 38 {
					name = name[:35] + "..."
				}
				fmt.Printf("%-6d %s %-10s %-40s %s\n", p.Day, statusIcon, p.Status, name, deps)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&statusFilter, "status", "s", "", "filtruj po statusie (todo, in_progress, done, blocked)")
	return cmd
}

func phaseShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show [day]",
		Short: "Szczegóły fazy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			day, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("day must be a number: %w", err)
			}

			store, err := openStore()
			if err != nil {
				return err
			}
			defer store.Close()

			p, err := store.GetPhaseByDay(day)
			if err != nil {
				return err
			}

			fmt.Printf("Dzień:       %d\n", p.Day)
			fmt.Printf("Nazwa:       %s\n", p.Name)
			fmt.Printf("Status:      %s %s\n", statusToIcon(p.Status), p.Status)
			if p.Description != "" {
				fmt.Printf("Opis:        %s\n", p.Description)
			}
			if len(p.DependsOn) > 0 {
				depStrs := make([]string, len(p.DependsOn))
				for i, d := range p.DependsOn {
					depStrs[i] = fmt.Sprintf("%d", d)
				}
				fmt.Printf("Zależności:  %s\n", strings.Join(depStrs, ", "))
			}
			fmt.Printf("Utworzono:    %s\n", p.CreatedAt.Format("2006-01-02 15:04"))
			fmt.Printf("Zaktualizowano: %s\n", p.UpdatedAt.Format("2006-01-02 15:04"))

			tasks, err := store.ListTasks(day)
			if err != nil {
				return err
			}
			if len(tasks) > 0 {
				fmt.Printf("\nZadania (%d):\n", len(tasks))
				for _, t := range tasks {
					check := "[ ]"
					if t.Done {
						check = "[x]"
					}
					fmt.Printf("  %s #%d %s\n", check, t.ID, t.Name)
				}
			}

			return nil
		},
	}
}

func phaseUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update [day] [status]",
		Short: "Zmień status fazy (todo, in_progress, done, blocked)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			day, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("day must be a number: %w", err)
			}
			status := model.PhaseStatus(args[1])

			store, err := openStore()
			if err != nil {
				return err
			}
			defer store.Close()

			if err := store.UpdatePhaseStatus(day, status); err != nil {
				return err
			}

			fmt.Printf("Faza day=%d → %s %s\n", day, statusToIcon(status), status)
			return nil
		},
	}
}

func phaseAddCmd() *cobra.Command {
	var day int
	var name, description string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Dodaj nową fazę",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := openStore()
			if err != nil {
				return err
			}
			defer store.Close()

			p := model.Phase{
				Day:         day,
				Name:        name,
				Description: description,
				Status:      model.StatusTodo,
			}

			id, err := store.AddPhase(p)
			if err != nil {
				return err
			}

			fmt.Printf("Dodano fazę #%d (day=%d): %s\n", id, day, name)
			return nil
		},
	}
	cmd.Flags().IntVarP(&day, "day", "d", 0, "numer dnia (wymagany)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "nazwa fazy (wymagana)")
	cmd.Flags().StringVar(&description, "desc", "", "opis fazy")
	cmd.MarkFlagRequired("day")
	cmd.MarkFlagRequired("name")
	return cmd
}

func taskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Zarządzanie zadaniami",
	}
	cmd.AddCommand(taskAddCmd())
	cmd.AddCommand(taskListCmd())
	cmd.AddCommand(taskDoneCmd())
	return cmd
}

func taskAddCmd() *cobra.Command {
	var phaseDay int
	var name string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Dodaj zadanie do fazy",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := openStore()
			if err != nil {
				return err
			}
			defer store.Close()

			if err := store.AddTask(phaseDay, name); err != nil {
				return err
			}

			fmt.Printf("Dodano zadanie '%s' do fazy day=%d\n", name, phaseDay)
			return nil
		},
	}
	cmd.Flags().IntVarP(&phaseDay, "phase", "p", 0, "numer dnia fazy (wymagany)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "nazwa zadania (wymagana)")
	cmd.MarkFlagRequired("phase")
	cmd.MarkFlagRequired("name")
	return cmd
}

func taskListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list [day]",
		Short: "Lista zadań fazy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			day, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("day must be a number: %w", err)
			}

			store, err := openStore()
			if err != nil {
				return err
			}
			defer store.Close()

			tasks, err := store.ListTasks(day)
			if err != nil {
				return err
			}

			if len(tasks) == 0 {
				fmt.Printf("Brak zadań dla fazy day=%d\n", day)
				return nil
			}

			fmt.Printf("Zadania fazy day=%d:\n", day)
			for _, t := range tasks {
				check := "[ ]"
				if t.Done {
					check = "[x]"
				}
				fmt.Printf("  %s #%d %s\n", check, t.ID, t.Name)
			}
			return nil
		},
	}
}

func taskDoneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "done [id]",
		Short: "Oznacz zadanie jako ukończone",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("task id must be a number: %w", err)
			}

			store, err := openStore()
			if err != nil {
				return err
			}
			defer store.Close()

			if err := store.CompleteTask(id); err != nil {
				return err
			}

			fmt.Printf("Zadanie #%d oznaczone jako ukończone\n", id)
			return nil
		},
	}
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Podsumowanie projektu",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := openStore()
			if err != nil {
				return err
			}
			defer store.Close()

			name, err := store.GetProjectName()
			if err != nil {
				return err
			}

			stats, err := store.GetStats()
			if err != nil {
				return err
			}

			fmt.Printf("Projekt: %s\n", name)
			fmt.Println(strings.Repeat("─", 40))
			fmt.Printf("Fazy:         %d\n", stats.TotalPhases)
			fmt.Printf("  TODO:       %d\n", stats.TodoPhases)
			fmt.Printf("  W toku:     %d\n", stats.InProgressPhases)
			fmt.Printf("  Gotowe:     %d\n", stats.DonePhases)
			fmt.Printf("  Zablokowane:%d\n", stats.BlockedPhases)
			fmt.Println(strings.Repeat("─", 40))
			fmt.Printf("Zadania:      %d (ukończone: %d)\n", stats.TotalTasks, stats.DoneTasks)
			fmt.Println(strings.Repeat("─", 40))

			barLen := 30
			filled := int(stats.PercentComplete / 100 * float64(barLen))
			bar := strings.Repeat("█", filled) + strings.Repeat("░", barLen-filled)
			fmt.Printf("Postęp:       [%s] %.1f%%\n", bar, stats.PercentComplete)

			return nil
		},
	}
}

func loadConfig() (*config.Config, error) {
	cfgPath := configFile
	if cfgPath == "" {
		cfgPath = config.FindConfigFile()
		if cfgPath == "" {
			cfgPath = config.DefaultConfigFile
		}
	}
	return config.Load(cfgPath)
}

func newLogger(cfg *config.Config) *logger.Logger {
	level := cfg.Log.Level
	if logLevel != "" {
		level = logLevel
	}
	format := cfg.Log.Format
	if logFormat != "" {
		format = logFormat
	}
	return logger.New(
		os.Stderr,
		logger.ParseLevel(level),
		logger.ParseFormat(format),
	)
}

func validatorOptions(cfg *config.Config) *validator.Options {
	opts := &validator.Options{PromoteContextFor: map[string]bool{}, RequiredDeps: map[string][]string{}, SectionAliases: cfg.SectionAliases}
	for _, dt := range cfg.Dependency.PromoteContextFor {
		opts.PromoteContextFor[dt] = true
	}
	if cfg.Dependency.FamilyRulesPath != "" {
		if fr, err := rules.LoadFamilyRules(cfg.Dependency.FamilyRulesPath); err == nil {
			for _, f := range fr.Families {
				if len(f.RequiredDeps) > 0 {
					opts.RequiredDeps[f.DocType] = append(opts.RequiredDeps[f.DocType], f.RequiredDeps...)
				}
			}
		}
	}
	return opts
}

func scanCmd() *cobra.Command {
	var indexOutput string
	var noTimestamps bool
	var deterministic bool
	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Skanuj repozytorium i wypisz pliki Markdown",
		Long:  "Przechodzi po katalogach zgodnie z konfiguracją, parsuje metadane i generuje indeks dokumentów.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}

			root := cfg.DocsRoot
			if len(args) > 0 {
				root = args[0]
			}

			log := newLogger(cfg)
			log.Info("Rozpoczynam skanowanie: root=%s", root)
			log.Debug("Ignorowane wzorce: %v", cfg.IgnorePatterns)

			idx, err := index.BuildIndex(root, cfg.IgnorePatterns)
			if err != nil {
				return fmt.Errorf("błąd skanowania: %w", err)
			}

			log.Info("Znaleziono %d dokumentów", idx.Count)

			fmt.Printf("%-50s %-20s %-15s %s\n", "ŚCIEŻKA", "DOC_ID", "STATUS", "TYP")
			fmt.Println(strings.Repeat("─", 100))
			for _, rec := range idx.All() {
				path := rec.Path
				if len(path) > 48 {
					path = "..." + path[len(path)-45:]
				}
				fmt.Printf("%-50s %-20s %-15s %s\n", path, rec.DocID, rec.Status, rec.DocType)
			}
			fmt.Println(strings.Repeat("─", 100))
			fmt.Printf("Razem: %d dokumentów\n", idx.Count)

			if indexOutput != "" {
				if deterministic {
					noTimestamps = true
				}
				if err := idx.SaveJSONWithOptions(indexOutput, index.SaveOptions{NoTimestamps: noTimestamps}); err != nil {
					return fmt.Errorf("błąd zapisu indeksu: %w", err)
				}
				log.Info("Indeks zapisany do %s", indexOutput)
				fmt.Printf("Indeks zapisany do: %s\n", indexOutput)
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&indexOutput, "output", "o", "", "zapisz indeks JSON do pliku (np. .docflow/cache/doc_index.json)")
	cmd.Flags().BoolVar(&noTimestamps, "no-timestamps", false, "zapisz indeks bez zmiennego created_at (deterministyczny output)")
	cmd.Flags().BoolVar(&deterministic, "deterministic", false, "alias dla --no-timestamps")
	return cmd
}

func validateCmd() *cobra.Command {
	var strict bool
	var warnOnly bool
	var oldIndexPath string
	var autoBump bool
	var saveIndex bool
	var statusAware bool
	var governanceRules string
	var noTimestamps bool
	var deterministic bool
	var format string
	var output string
	var againstPath string
	var failOn string
	var show string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Waliduj metadane i strukturę dokumentów",
		Long:  "Sprawdza frontmatter, wymagane pola, unikalność doc_id oraz zgodność nazw plików.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if strict && warnOnly {
				return commandUsageError(
					format,
					output,
					"DOCFLOW.CLI.CONFLICTING_FLAGS",
					"--strict i --warn nie mogą być użyte jednocześnie",
					map[string]any{"flags": []string{"--strict", "--warn"}},
				)
			}
			if format != "text" && format != "json" && format != "sarif" {
				return commandUsageError(
					format,
					output,
					"DOCFLOW.CLI.INVALID_FORMAT",
					fmt.Sprintf("nieobsługiwany format: %s (dozwolone: text|json|sarif)", format),
					map[string]any{"format": format},
				)
			}
			if failOn != "any" && failOn != "new" {
				return commandUsageError(
					format,
					output,
					"DOCFLOW.CLI.INVALID_FAIL_ON",
					fmt.Sprintf("nieobsługiwane --fail-on: %s (dozwolone: any|new)", failOn),
					map[string]any{"fail_on": failOn},
				)
			}
			if show != "all" && show != "new" && show != "existing" {
				return commandUsageError(
					format,
					output,
					"DOCFLOW.CLI.INVALID_SHOW",
					fmt.Sprintf("nieobsługiwane --show: %s (dozwolone: all|new|existing)", show),
					map[string]any{"show": show},
				)
			}
			if failOn == "new" && againstPath == "" {
				return commandUsageError(
					format,
					output,
					"DOCFLOW.CLI.MISSING_AGAINST",
					"podaj --against dla --fail-on new",
					map[string]any{"flag": "--against"},
				)
			}
			if (format == "json" || format == "sarif") && output == "" {
				output = "-"
			}

			cfg, err := loadConfig()
			if err != nil {
				return commandRuntimeError(
					format,
					output,
					"DOCFLOW.RUNTIME.CONFIG_LOAD_FAILED",
					"nie można wczytać konfiguracji",
					err,
					nil,
				)
			}

			log := newLogger(cfg)
			log.Info("Walidacja dokumentów: root=%s", cfg.DocsRoot)

			opts := validatorOptions(cfg)
			opts.StatusAware = statusAware
			if governanceRules != "" {
				gr, err := governance.Load(governanceRules)
				if err != nil {
					return commandRuntimeError(
						format,
						output,
						"DOCFLOW.RUNTIME.GOVERNANCE_LOAD_FAILED",
						fmt.Sprintf("nie można wczytać reguł governance z %s", governanceRules),
						err,
						map[string]any{"path": governanceRules},
					)
				}
				opts.GovernanceRules = gr
			}
			execRes, err := collectExecutionResult("validate", cfg)
			if err != nil {
				return commandRuntimeError(
					format,
					output,
					"DOCFLOW.RUNTIME.EXECUTION_LAYER_FAILED",
					"nie można przygotować execution layer",
					err,
					nil,
				)
			}
			report, err := validator.ValidateFacts(execRes.Facts, execRes.Index, opts)
			if err != nil {
				return commandRuntimeError(
					format,
					output,
					"DOCFLOW.RUNTIME.VALIDATE_FAILED",
					"nie można wykonać walidacji dokumentów",
					err,
					nil,
				)
			}

			if format == "text" {
				if len(report.Issues) == 0 {
					fmt.Println("✅ Brak problemów z metadanymi.")
				} else {
					fmt.Printf("%s\n", strings.Repeat("─", 100))
					fmt.Printf("%-6s %-12s %-45s %-30s\n", "LEVEL", "TYPE", "FILE", "MESSAGE")
					fmt.Printf("%s\n", strings.Repeat("─", 100))
					for _, issue := range report.SortedIssues() {
						file := issue.File
						if len(file) > 44 {
							file = "..." + file[len(file)-41:]
						}
						msg := issue.Message
						if issue.DocID != "" {
							msg = fmt.Sprintf("%s (doc_id=%s)", msg, issue.DocID)
						}
						fmt.Printf("%-6s %-12s %-45s %-30s\n", strings.ToUpper(string(issue.Level)), issue.Type, file, msg)
					}
					fmt.Printf("%s\n", strings.Repeat("─", 100))
					fmt.Printf("Pliki: %d | Dokumenty: %d | Błędy: %d | Ostrzeżenia: %d\n", report.Files, report.Documents, report.ErrorCount(), report.WarnCount())
				}
			}

			// Change tracking vs cached index.
			var oldPath string
			if oldIndexPath != "" {
				oldPath = oldIndexPath
			} else {
				def := filepath.Join(cfg.Cache.Dir, "doc_index.json")
				if _, err := os.Stat(def); err == nil {
					oldPath = def
				}
			}

			var idx *index.DocumentIndex
			idx = execRes.Index
			if oldPath != "" || saveIndex || autoBump {
				if idx == nil {
					idx, err = index.BuildIndex(cfg.DocsRoot, cfg.IgnorePatterns)
					if err != nil {
						return commandRuntimeError(
							format,
							output,
							"DOCFLOW.RUNTIME.INDEX_BUILD_FAILED",
							"nie można zbudować indeksu dokumentów",
							err,
							map[string]any{"docs_root": cfg.DocsRoot},
						)
					}
				}
			}

			if oldPath != "" && idx != nil {
				oldIdx, err := index.LoadJSON(oldPath)
				if err != nil {
					log.Warn("Nie można załadować starego indeksu %s: %v", oldPath, err)
				} else {
					changes := index.DocDiff(oldIdx, idx)
					if format == "text" {
						if len(changes) == 0 {
							fmt.Printf("Brak zmian względem %s\n", oldPath)
						} else {
							fmt.Printf("\nZmiany względem %s:\n", oldPath)
							fmt.Printf("%-15s %-20s %s\n", "TYPE", "DOC_ID", "PATH")
							for _, c := range changes {
								fmt.Printf("%-15s %-20s %s\n", c.Type, c.DocID, c.Path)
							}
						}
					}
					if autoBump && len(changes) > 0 {
						bumps, err := index.AutoBumpVersions(changes, cfg.DocsRoot)
						if err != nil {
							return commandRuntimeError(
								format,
								output,
								"DOCFLOW.RUNTIME.AUTO_BUMP_FAILED",
								"nie można wykonać auto-bump wersji",
								err,
								nil,
							)
						}
						if len(bumps) > 0 && format == "text" {
							fmt.Printf("\nAuto-bump wersji:\n")
							for _, b := range bumps {
								fmt.Printf("%-50s %s -> %s (%s)\n", b.Path, b.From, b.To, b.Reason)
							}
						}
						if len(bumps) > 0 && saveIndex {
							idx, err = index.BuildIndex(cfg.DocsRoot, cfg.IgnorePatterns)
							if err != nil {
								return commandRuntimeError(
									format,
									output,
									"DOCFLOW.RUNTIME.INDEX_BUILD_FAILED",
									"nie można zbudować indeksu dokumentów po auto-bump",
									err,
									map[string]any{"docs_root": cfg.DocsRoot},
								)
							}
						}
					}
				}
			}

			if saveIndex {
				if idx == nil {
					idx, err = index.BuildIndex(cfg.DocsRoot, cfg.IgnorePatterns)
					if err != nil {
						return commandRuntimeError(
							format,
							output,
							"DOCFLOW.RUNTIME.INDEX_BUILD_FAILED",
							"nie można zbudować indeksu dokumentów do zapisu",
							err,
							map[string]any{"docs_root": cfg.DocsRoot},
						)
					}
				}
				path := oldIndexPath
				if path == "" {
					path = filepath.Join(cfg.Cache.Dir, "doc_index.json")
				}
				if deterministic {
					noTimestamps = true
				}
				if err := idx.SaveJSONWithOptions(path, index.SaveOptions{NoTimestamps: noTimestamps}); err != nil {
					return commandRuntimeError(
						format,
						output,
						"DOCFLOW.RUNTIME.OUTPUT_WRITE_FAILED",
						fmt.Sprintf("nie można zapisać indeksu do %s", path),
						err,
						map[string]any{"output": path},
					)
				}
				log.Info("Indeks zapisany do %s", path)
			}

			reportJSON := buildValidateJSONReport(report)
			reportJSON.FailOn = failOn
			reportJSON.Show = show
			if againstPath != "" {
				reportJSON, err = applyValidateBaseline(reportJSON, againstPath, show)
				if err != nil {
					return commandRuntimeError(
						format,
						output,
						"DOCFLOW.RUNTIME.BASELINE_LOAD_FAILED",
						fmt.Sprintf("nie można wczytać baseline validate z %s", againstPath),
						err,
						map[string]any{"path": againstPath},
					)
				}
			}

			if format == "json" {
				if err := writeValidateJSONReport(reportJSON, output); err != nil {
					return commandRuntimeError(
						format,
						output,
						"DOCFLOW.RUNTIME.OUTPUT_WRITE_FAILED",
						fmt.Sprintf("nie można zapisać raportu validate do %s", output),
						err,
						map[string]any{"output": output},
					)
				}
			} else if format == "sarif" {
				if err := writeValidateSARIFReport(reportJSON, output); err != nil {
					return commandRuntimeError(
						format,
						output,
						"DOCFLOW.RUNTIME.OUTPUT_WRITE_FAILED",
						fmt.Sprintf("nie można zapisać raportu validate SARIF do %s", output),
						err,
						map[string]any{"output": output},
					)
				}
			}

			errorsForGate := report.ErrorCount()
			if failOn == "new" {
				errorsForGate = reportJSON.NewErrorCount
			}

			if strict && errorsForGate > 0 && !warnOnly {
				if failOn == "new" {
					return fmt.Errorf("walidacja zakończona: %d nowych błędów", errorsForGate)
				}
				return fmt.Errorf("walidacja zakończona: %d błędów", errorsForGate)
			}
			if warnOnly {
				if failOn == "new" {
					log.Warn("Walidacja uruchomiona w trybie warn — nowe błędy nie blokują (new=%d)", reportJSON.NewErrorCount)
				} else {
					log.Warn("Walidacja uruchomiona w trybie warn — błędy nie blokują")
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&strict, "strict", false, "zwróć błąd exit gdy wykryto błędy")
	cmd.Flags().BoolVar(&warnOnly, "warn", false, "raportuj, ale nie blokuj nawet przy błędach")
	cmd.Flags().StringVar(&oldIndexPath, "old-index", "", "ścieżka do poprzedniego indeksu (domyślnie .docflow/cache/doc_index.json jeśli istnieje)")
	cmd.Flags().BoolVar(&autoBump, "auto-bump", false, "automatycznie podbij wersję zmienionych dokumentów")
	cmd.Flags().BoolVar(&saveIndex, "save-index", false, "zapisz bieżący indeks do cache (po walidacji)")
	cmd.Flags().BoolVar(&noTimestamps, "no-timestamps", false, "zapisz indeks bez zmiennego created_at (wymaga --save-index)")
	cmd.Flags().BoolVar(&deterministic, "deterministic", false, "alias dla --no-timestamps")
	cmd.Flags().BoolVar(&statusAware, "status-aware", false, "włącz progressive validation zależnie od statusu dokumentu")
	cmd.Flags().StringVar(&governanceRules, "governance", "", "ścieżka do GOVERNANCE_RULES.yaml (opcjonalne)")
	cmd.Flags().StringVar(&format, "format", "text", "format: text|json|sarif")
	cmd.Flags().StringVar(&output, "output", "", "plik wyjściowy (json/sarif); '-' oznacza stdout")
	cmd.Flags().StringVar(&againstPath, "against", "", "ścieżka do baseline validate JSON")
	cmd.Flags().StringVar(&failOn, "fail-on", "any", "gating: any|new")
	cmd.Flags().StringVar(&show, "show", "all", "które issues zwrócić w JSON/SARIF: all|new|existing")
	return cmd
}

type validateIssueJSON struct {
	Code    string         `json:"code"`
	Level   string         `json:"level"`
	Type    string         `json:"type"`
	Path    string         `json:"path"`
	DocID   string         `json:"doc_id,omitempty"`
	Line    int            `json:"line,omitempty"`
	Column  int            `json:"column,omitempty"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

type validateReportJSON struct {
	SchemaVersion      string              `json:"schema_version"`
	IdentityVersion    string              `json:"identity_version,omitempty"`
	Baseline           *baselineMetaJSON   `json:"baseline,omitempty"`
	FailOn             string              `json:"fail_on,omitempty"`
	Show               string              `json:"show,omitempty"`
	Files              int                 `json:"files"`
	Documents          int                 `json:"documents"`
	ErrorCount         int                 `json:"error_count"`
	WarnCount          int                 `json:"warn_count"`
	NewErrorCount      int                 `json:"new_error_count,omitempty"`
	NewWarnCount       int                 `json:"new_warn_count,omitempty"`
	ExistingErrorCount int                 `json:"existing_error_count,omitempty"`
	ExistingWarnCount  int                 `json:"existing_warn_count,omitempty"`
	Issues             []validateIssueJSON `json:"issues"`
}

func buildValidateJSONReport(report *validator.Report) validateReportJSON {
	canonical := engine.BuildValidateReport(report)
	issues := make([]validateIssueJSON, 0, len(canonical.Issues))
	for _, issue := range canonical.Issues {
		line := 0
		column := 0
		if issue.Location != nil {
			line = issue.Location.Line
			column = issue.Location.Column
		}
		issues = append(issues, validateIssueJSON{
			Code:    issue.Code,
			Level:   string(issue.Severity),
			Type:    issue.Type,
			Path:    issue.Path,
			DocID:   issue.DocID,
			Line:    line,
			Column:  column,
			Message: issue.Message,
			Details: issue.Details,
		})
	}

	return validateReportJSON{
		SchemaVersion:   canonical.SchemaVersion,
		IdentityVersion: currentIdentityVersion(),
		Files:           canonical.Files,
		Documents:       canonical.Documents,
		ErrorCount:      canonical.ErrorCount,
		WarnCount:       canonical.WarnCount,
		Issues:          issues,
	}
}

func writeValidateJSONReport(report validateReportJSON, output string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("błąd serializacji raportu validate: %w", err)
	}
	if output == "" || output == "-" {
		_, err = os.Stdout.Write(append(data, '\n'))
		return err
	}
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return fmt.Errorf("błąd utworzenia katalogu dla raportu validate: %w", err)
	}
	return os.WriteFile(output, data, 0o644)
}

func planCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan",
		Short: "Pokaż plan generowania dokumentów",
		Long:  "Analizuje graf zależności i wyświetla kolejność generowania dokumentów.",
	}
	cmd.AddCommand(planDailyCmd())
	return cmd
}

func analyzePatternsCmd() *cobra.Command {
	var topN int
	var schemaDir string
	cmd := &cobra.Command{
		Use:   "analyze-patterns",
		Short: "Analizuj wzorce sekcji w szablonach",
		Long:  "Skanuje szablony Markdown, grupuje sekwencje nagłówków, wylicza top wzorce i generuje schematy sekcji.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}

			log := newLogger(cfg)
			log.Info("Analiza wzorców sekcji: root=%s", cfg.DocsRoot)

			patterns, err := pattern.ScanPatterns(cfg.DocsRoot, cfg.IgnorePatterns)
			if err != nil {
				return fmt.Errorf("błąd skanowania wzorców: %w", err)
			}

			report := pattern.GroupPatterns(patterns)
			log.Info("Znaleziono %d szablonów, %d unikalnych wzorców", report.TotalTemplates, report.UniquePatterns)

			fmt.Print(pattern.FormatReport(report, topN))

			if schemaDir != "" {
				schemas := pattern.GenerateSchemas(report, topN)
				if err := pattern.SaveAllSchemas(schemas, schemaDir); err != nil {
					return fmt.Errorf("błąd zapisu schematów: %w", err)
				}
				fmt.Printf("\nZapisano %d schematów do: %s\n", len(schemas), schemaDir)
			}

			return nil
		},
	}
	cmd.Flags().IntVarP(&topN, "top", "n", 10, "liczba top wzorców do wyświetlenia")
	cmd.Flags().StringVarP(&schemaDir, "schema-dir", "s", "", "katalog do zapisu section_schema.yaml (np. .docflow/schemas)")
	return cmd
}

func statusToIcon(s model.PhaseStatus) string {
	switch s {
	case model.StatusTodo:
		return "○"
	case model.StatusInProgress:
		return "◐"
	case model.StatusDone:
		return "●"
	case model.StatusBlocked:
		return "✖"
	default:
		return "?"
	}
}
