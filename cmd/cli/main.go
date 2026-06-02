package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/yellowkode-academy/linkvault/internal/link"
	"github.com/yellowkode-academy/linkvault/internal/storage"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func openRepo() (*storage.SQLiteRepository, error) {
	dsn := os.Getenv("LINKVAULT_DB")
	if dsn == "" {
		dsn = "linkvault.db"
	}
	return storage.NewSQLiteRepository(dsn)
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "linkvault",
		Short:         "Gerenciador de links pessoal",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(newAddCmd(), newListCmd(), newSearchCmd())
	return cmd
}

func newAddCmd() *cobra.Command {
	var title, tags string
	cmd := &cobra.Command{
		Use:   "add <url>",
		Short: "Salva um link novo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := openRepo()
			if err != nil {
				return err
			}
			defer repo.Close()
			l := link.NewLink(args[0], title, tags)
			if err := l.Validate(); err != nil {
				return err
			}
			saved, err := repo.Save(context.Background(), l)
			if err != nil {
				return err
			}
			fmt.Printf("Link salvo com ID %d\n", saved.ID)
			return nil
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "Titulo do link")
	cmd.Flags().StringVar(&tags, "tags", "", "Tags separadas por virgula")
	cmd.MarkFlagRequired("title")
	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lista todos os links",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := openRepo()
			if err != nil {
				return err
			}
			defer repo.Close()
			links, err := repo.List(context.Background())
			if err != nil {
				return err
			}
			printLinks(links)
			return nil
		},
	}
}

func newSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Busca links por texto",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := openRepo()
			if err != nil {
				return err
			}
			defer repo.Close()
			links, err := repo.Search(context.Background(), args[0])
			if err != nil {
				return err
			}
			printLinks(links)
			return nil
		},
	}
}

func printLinks(links []link.Link) {
	if len(links) == 0 {
		fmt.Println("Nenhum link encontrado.")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tURL\tTITULO\tTAGS")
	for _, l := range links {
		tags := strings.Replace(l.Tags, ",", " ", -1)
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", l.ID, l.URL, l.Title, tags)
	}
	w.Flush()
}
