package store

// QueryMaps runs a query and returns each row as column → value.
func (db *DB) QueryMaps(query string, args ...any) (cols []string, out []map[string]any, err error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	if cols, err = rows.Columns(); err != nil {
		return nil, nil, err
	}
	out = []map[string]any{}
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, nil, err
		}
		m := make(map[string]any, len(cols))
		for i, c := range cols {
			if b, ok := vals[i].([]byte); ok {
				vals[i] = string(b)
			}
			m[c] = vals[i]
		}
		out = append(out, m)
	}
	return cols, out, rows.Err()
}
