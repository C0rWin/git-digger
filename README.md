# git-digger

## Overview

`git-digger` is a tool designed to analyze git repositories and generate interaction graphs between project maintainers and the git projects they contribute to. This tool aims to provide insights into the knowledge and expertise of code contributors, making it easier to understand the dynamics of software development within your organization or open-source projects.

The generated graph is compatible with Gephi and is based on the GraphML format.

## Features

- Scan a folder containing multiple git repositories
- Generate a Gephi-compatible interaction graph
- Filter contributions based on a specific date range

## Installation

```bash
git clone https://github.com/yourusername/git-digger.git
cd git-digger
make install
```

## Usage

To get started with `git-digger`, you can run the following command:

```bash
git-digger [flags]
```

### Flags

- `-h, --help`: Show help for `git-digger`
- `-r, --repo string`: Path to the git repository or folder containing multiple repositories
- `-s, --since string`: Filter contributions since this date (default is "02/01/2006")

### Example

```bash
git-digger --repo /path/to/repositories --since "01/01/2021"
```

## Output Format

The output is a GraphML file with the following structure:

```go
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
```

## Contributing

We welcome contributions! Please see the [CONTRIBUTING.md](CONTRIBUTING.md) for more details.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Support

For any questions or issues, please open an issue on GitHub or contact the maintainers.

---

Happy digging! 🛠️