// Package sqloader provides a way to read named SQL queries from SQL files.
package sqloader

import (
	"bufio"
	"errors"
	"io"
	"os"
)

type SQLoader struct {
	queries map[string]string
}

// NewSQLoader loads a SQL file and returns a pointer to SQLoader with loaded queries.
func NewSQLoader(sqlfile string) (*SQLoader, error) {
	sqloaderPtr := &SQLoader{queries: map[string]string{}}
	var err error

	f, err := os.Open(sqlfile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	var queryName, query string
	var r rune
	var extractingQuery bool
	for {
		r, _, err = reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if !extractingQuery {
			// If not extracting a query, we should scan for a "-- name: queryName"

			// Match the first '-', otherwise every rune that is not in a query is discarded
			if r == '-' {
				// Read the second '-'
				r, _, err = reader.ReadRune()
				if err != nil {
					if err == io.EOF {
						return nil, errors.New("Malformed SQLoader file")
					}
					return nil, err
				}

				// We expect another '-'
				if r != '-' {
					return nil, errors.New("Malformed SQLoader file")
				}

				// Discard all white spaces
				for {
					r, _, err = reader.ReadRune()
					if err != nil {
						if err == io.EOF {
							return nil, errors.New("Malformed SQLoader file")
						}
						return nil, err
					}

					if r != ' ' {
						break
					}
				}
				// Unread one rune, that breaked the discarding white spaces loop
				err = reader.UnreadRune()
				if err != nil {
					return nil, err
				}

				// If the next rune is '/', it starts a query name
				r, _, err = reader.ReadRune()
				if err != nil {
					if err == io.EOF {
						return nil, errors.New("Malformed SQLoader file")
					}
					return nil, err
				}
				// If the rune is not '/', it is just a normal SQL comment
				if r != '/' {
					continue
				}

				// Read the name until match a new line
				queryName = ""
				for {
					r, _, err = reader.ReadRune()
					if err != nil {
						if err == io.EOF {
							return nil, errors.New("Malformed SQLoader file")
						}
						return nil, err
					}

					// CRLF
					if r == '\r' {
						r, _, err = reader.ReadRune()
						if err != nil {
							if err == io.EOF {
								return nil, errors.New("Malformed SQLoader file")
							}
							return nil, err
						}

						if r != '\n' {
							return nil, errors.New("Malformed SQLoader file")
						}

						// New line found
						break
					}

					// LF
					if r == '\n' {
						break
					}

					queryName += string(r)
				}
				if queryName == "" {
					return nil, errors.New("Malformed SQLoader file, missing name for SQL query")
				}

				extractingQuery = true
				query = ""
			}
		} else {
			// We are extracting a query, read until "-- end"

			// First check if the query is being ended
			if r == '-' {
				// Read the second '-'
				r, _, err = reader.ReadRune()
				if err != nil {
					if err == io.EOF {
						return nil, errors.New("Malformed SQLoader file")
					}
					return nil, err
				}

				// We expect another '-'
				if r != '-' {
					return nil, errors.New("Malformed SQLoader file")
				}

				// Discard all white spaces
				for {
					r, _, err = reader.ReadRune()
					if err != nil {
						if err == io.EOF {
							return nil, errors.New("Malformed SQLoader file")
						}
						return nil, err
					}

					if r != ' ' {
						break
					}
				}
				// Unread one rune, that breaked the discarding white spaces loop
				err = reader.UnreadRune()
				if err != nil {
					return nil, err
				}

				// If the next rune match '/', the query ended
				r, _, err = reader.ReadRune()
				if err != nil {
					if err == io.EOF {
						return nil, errors.New("Malformed SQLoader file")
					}
					return nil, err
				}
				// If the rune is not '/', it is just a normal SQL comment
				if r != '/' {
					continue
				}

				extractingQuery = false

				sqloaderPtr.queries[queryName] = query

				continue
			}

			query += string(r)
		}
	}

	return sqloaderPtr, nil
}

// Get gets a named query and returns it. If no query is found, it returns an empty string.
func (sqloader *SQLoader) Get(queryName string) string {
	query, ok := sqloader.queries[queryName]
	if !ok {
		return ""
	}
	return query
}
