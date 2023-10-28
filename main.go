package main

import (
	"fmt"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/spf13/cobra"
)

// Contributors is a list of contributors
type Contributors map[string]struct{}

// Graph is a map of contributors
type Graph map[string]Contributors

// Repository is a git repository
type Repository struct {
	Path  string
	Graph Graph
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

	h, err := g.Head()
	if err != nil {
		return nil
	}

	// now := time.Now()
	// since := now.AddDate(0, -1, 0)
	// l, err := g.Log(&git.LogOptions{All: true, Since: &since})
	l, err := g.Log(&git.LogOptions{From: h.Hash()})
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
		repoPath  string
		graphType string
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
				Graph: make(Graph),
			}

			filepath.Walk(repoPath, r.Scan)

			g := graphviz.New()
			graph, err := g.Graph()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			// keep track of contributors to avoid duplicates
			var contributorNodes map[string]*cgraph.Node

			for name, repo := range r.Graph {
				rV, err := graph.CreateNode(name)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				for contributor := range repo {
					cV, exists := contributorNodes[contributor]
					if !exists {
						cV, err = graph.CreateNode(contributor)
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
					}
					graph.CreateEdge("", cV, rV)
				}
			}

			graph.SetLayout(graphType)

			if err := g.RenderFilename(graph, graphviz.XDOT, "graph.gv"); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&repoPath, "repo", "r", "", "Path to the git repository")
	cmd.Flags().StringVarP(&graphType, "type", "t", "circo", "Type of the graph")
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
