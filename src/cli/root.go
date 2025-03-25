package cli

import (
	"os"

	"tagstore/src/config"
	"tagstore/src/lib"
	"tagstore/src/repository"
	"tagstore/src/service"

	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command
	cfg     *config.Config
	db      *lib.DB
	tagRepo *repository.TagRepository
	fileRepo *repository.FileRepository
	fileTagRepo *repository.FileTagRepository
	tagService *service.TagService
	searchService *service.SearchService
)

// Execute runs the root command
func Execute() {
	var err error

	cfg, err = config.Load()
	if err != nil {
		lib.Get().Error("Failed to load config: %v", err)
		os.Exit(1)
	}

	lib.Init(cfg.LogLevel)

	db, err = lib.Connect(cfg.DatabasePath)
	if err != nil {
		lib.Get().Error("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	tagRepo = repository.NewTagRepository(db)
	fileRepo = repository.NewFileRepository(db)
	fileTagRepo = repository.NewFileTagRepository(db)
	tagService = service.NewTagService(tagRepo, fileRepo, fileTagRepo, lib.Get())
	searchService = service.NewSearchService(fileTagRepo, lib.Get())

	rootCmd = &cobra.Command{
		Use:   "tagstore",
		Short: "A CLI tool for tagging and searching files",
		Long:  "TagStore helps you organize files with arbitrary tags and quickly locate them using boolean tag queries.",
	}

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(importCmd)

	if err := rootCmd.Execute(); err != nil {
		lib.Get().Error("Command failed: %v", err)
		os.Exit(1)
	}
}
