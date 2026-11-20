package cli

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

var sqlCmd = &cobra.Command{
	Use:   "sql <query>",
	Short: "Run a read-only SQL query on the local database (JSON rows)",
	Long: `Run a read-only SQL query on the local database and print the rows as a
JSON array of objects, columns in query order.

Tables: forums, topics, posts, users, and posts_fts (FTS5 over posts.title
and posts.text). Times are UTC text like 2026-10-12T10:53:00Z. Post ids grow
over time, so ORDER BY id is oldest first. '.schema' prints the schema.`,
	Example: `  mz sql 'SELECT count(*) AS posts FROM posts'
  mz sql "SELECT u.name, count(*) n FROM posts p JOIN users u ON u.id = p.author_id GROUP BY 1 ORDER BY n DESC LIMIT 20"
  mz sql "SELECT p.id, snippet(posts_fts, 1, '[', ']', '…', 20) FROM posts_fts JOIN posts p ON p.id = posts_fts.rowid WHERE posts_fts MATCH 'training camp' LIMIT 5"`,
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := openDB(true)
		if err != nil {
			return err
		}
		defer db.Close()
		q := args[0]
		if q == ".schema" {
			q = `SELECT sql FROM sqlite_schema WHERE sql IS NOT NULL AND name NOT LIKE 'posts_fts_%'`
		}
		cols, rows, err := db.QueryMaps(q)
		if err != nil {
			return err
		}
		return printRows(cols, rows)
	},
}

// printRows writes rows as a JSON array, keeping the column order of the query.
func printRows(cols []string, rows []map[string]any) error {
	var b bytes.Buffer
	b.WriteString("[")
	for i, r := range rows {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString("\n  {")
		for j, c := range cols {
			if j > 0 {
				b.WriteString(", ")
			}
			if err := writeJSON(&b, c); err != nil {
				return err
			}
			b.WriteString(": ")
			if err := writeJSON(&b, r[c]); err != nil {
				return err
			}
		}
		b.WriteString("}")
	}
	if len(rows) > 0 {
		b.WriteString("\n")
	}
	b.WriteString("]\n")
	_, err := os.Stdout.Write(b.Bytes())
	return err
}

func writeJSON(b *bytes.Buffer, v any) error {
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return err
	}
	b.Truncate(b.Len() - 1) // Encode's newline
	return nil
}

func init() {
	rootCmd.AddCommand(sqlCmd)
}
