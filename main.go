package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

// Contributors is a list of contributors
type Contributors map[string]struct{}

// Graph is a map of contributors
type Graph map[string]Contributors

// Repository is a git repository
type Repository struct {
	Path  string
	Since string
	Graph Graph
}

// GraphML is a graphml file
type GraphML struct {
	XMLName xml.Name `xml:"graphml"`
	Nodes   []Node   `xml:"graph>node"`
	Edges   []Edge   `xml:"graph>edge"`
}

// Node is a node in the graphml file
type Node struct {
	ID   string `xml:"id,attr"`
	Type string `xml:"type,attr"`
}

// Edge is an edge in the graphml file
type Edge struct {
	Source string `xml:"source,attr"`
	Target string `xml:"target,attr"`
}

func (r *Repository) Scan(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err)
		return err
	}

	g, err := git.PlainOpen(path)
	if err != nil {
		return nil
	}

	since, err := time.Parse("02/01/2006", r.Since)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	l, err := g.Log(&git.LogOptions{All: true, Since: &since})
	if err != nil {
		return nil
	}

	r.Graph[path] = make(Contributors)
	err = l.ForEach(func(c *object.Commit) error {
		r.Graph[path][c.Author.Email] = struct{}{}
		return nil
	})
	return filepath.SkipDir
}

func main() {
	var (
		repoPath string
		since    string
	)

	cmd := cobra.Command{
		Use:   "git-digger",
		Short: "git-digger is a tool to dig into git repositories",
		Long:  "git-digger is a tool to dig into git repositories",
		Run: func(cmd *cobra.Command, args []string) {
			if len(repoPath) == 0 {
				cmd.Help()
				os.Exit(1)
			}

			r := Repository{
				Path:  repoPath,
				Since: since,
				Graph: make(Graph),
			}

			filepath.Walk(repoPath, r.Scan)

			graph := GraphML{}
			for repo, contributors := range r.Graph {
				if len(repo) == 0 {
					continue
				}
				graph.Nodes = append(graph.Nodes, Node{ID: repo, Type: "project"})
				for contributor := range contributors {
					if len(contributor) == 0 {
						continue
					}
					graph.Nodes = append(graph.Nodes, Node{ID: contributor, Type: "contributor"})
					graph.Edges = append(graph.Edges, Edge{Source: contributor, Target: repo})
				}
			}

			file, err := os.Create("graph.graphml")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer file.Close()

			enc := xml.NewEncoder(file)
			enc.Indent("", "  ")
			if err := enc.Encode(graph); err != nil {
				fmt.Printf("error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&repoPath, "repo", "r", "", "Path to the git repository")
	cmd.Flags().StringVarP(&since, "since", "s", "02/01/2006", "Since date")
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
