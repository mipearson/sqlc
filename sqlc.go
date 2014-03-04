// Package sqlc helps you combine arbitrary bits of SQL
package sqlc

import (
	"fmt"
	"strings"
)

type component struct {
	partial string
	args    []interface{}
}

/*
From the psql documentation:
SELECT [ ALL | DISTINCT [ ON ( expression [, ...] ) ] ]
    * | expression [ [ AS ] output_name ] [, ...]
    [ FROM from_item [, ...] ]
    [ WHERE condition ]
    [ GROUP BY expression [, ...] ]
    [ HAVING condition [, ...] ]
    [ WINDOW window_name AS ( window_definition ) [, ...] ]
    [ { UNION | INTERSECT | EXCEPT } [ ALL | DISTINCT ] select ]
    [ ORDER BY expression [ ASC | DESC | USING operator ] [ NULLS { FIRST | LAST } ] [, ...] ]
    [ LIMIT { count | ALL } ]
    [ OFFSET start [ ROW | ROWS ] ]
    [ FETCH { FIRST | NEXT } [ count ] { ROW | ROWS } ONLY ]
    [ FOR { UPDATE | NO KEY UPDATE | SHARE | KEY SHARE } [ OF table_name [, ...] ] [ NOWAIT ] [...] ]

From the mysql documentation:
SELECT
    [ALL | DISTINCT | DISTINCTROW ]
      [HIGH_PRIORITY]
      [STRAIGHT_JOIN]
      [SQL_SMALL_RESULT] [SQL_BIG_RESULT] [SQL_BUFFER_RESULT]
      [SQL_CACHE | SQL_NO_CACHE] [SQL_CALC_FOUND_ROWS]
    select_expr [, select_expr ...]
    [FROM table_references
      [PARTITION partition_list]
    [WHERE where_condition]
    [GROUP BY {col_name | expr | position}
      [ASC | DESC], ... [WITH ROLLUP]]
    [HAVING where_condition]
    [ORDER BY {col_name | expr | position}
      [ASC | DESC], ...]
    [LIMIT {[offset,] row_count | row_count OFFSET offset}]
    [PROCEDURE procedure_name(argument_list)]
    [INTO OUTFILE 'file_name'
        [CHARACTER SET charset_name]
        export_options
      | INTO DUMPFILE 'file_name'
      | INTO var_name [, var_name]]
    [FOR UPDATE | LOCK IN SHARE MODE]]
*/

// Statement is a SQL string being built
type Statement struct {
	// PostgreSQL replaces "?, ?" with "$1, $2" if true
	PostgreSQL bool
	selects    []component
	froms      []component
	joins      []component
	wheres     []component
	groups     []component
	havings    []component
	orders     []component
	limit      component
}

// Select adds a SELECT stanza, joined by commas
func (s Statement) Select(partial string, args ...interface{}) Statement {
	s.selects = append(s.selects, component{partial, args})
	return s
}

// From adds a FROM stanza, joined by commas
func (s Statement) From(partial string, args ...interface{}) Statement {
	s.froms = append(s.froms, component{partial, args})
	return s
}

// Join adds a JOIN stanza, joined by spaces
// Unlike other stanzas, you need to specify the JOIN/INNER JOIN/LEFT JOIN bit yourself.
func (s Statement) Join(partial string, args ...interface{}) Statement {
	s.joins = append(s.joins, component{partial, args})
	return s
}

// Where adds a WHERE stanza, wrapped in brackets and joined by AND
func (s Statement) Where(partial string, args ...interface{}) Statement {
	s.wheres = append(s.wheres, component{"(" + partial + ")", args})
	return s
}

// Having adds a HAVING stanza, wrapped in brackets and joined by AND
func (s Statement) Having(partial string, args ...interface{}) Statement {
	s.havings = append(s.havings, component{"(" + partial + ")", args})
	return s
}

// Group adds a GROUP BY stanza, joined by commas
func (s Statement) Group(partial string, args ...interface{}) Statement {
	s.groups = append(s.groups, component{partial, args})
	return s
}

// Order adds an ORDER BY stanza, joined by commas
func (s Statement) Order(partial string, args ...interface{}) Statement {
	s.orders = append(s.orders, component{partial, args})
	return s
}

// Limit sets or overwrites the LIMIT stanza
func (s Statement) Limit(partial string, args ...interface{}) Statement {
	s.limit = component{partial, args}
	return s
}

// Args returns positional arguments in the order they will appear in the SQL.
func (s Statement) Args() []interface{} {
	args := make([]interface{}, 0)

	appendArgs(&args, s.selects...)
	appendArgs(&args, s.froms...)
	appendArgs(&args, s.joins...)
	appendArgs(&args, s.wheres...)
	appendArgs(&args, s.groups...)
	appendArgs(&args, s.havings...)
	appendArgs(&args, s.orders...)
	appendArgs(&args, s.limit)

	return args
}

// SQL joins your stanzas, returning the composed SQL.
func (s Statement) SQL() string {
	parts := make([]string, 0)

	if len(s.selects) > 0 {
		parts = append(parts, "SELECT "+joinParts(s.selects, ", "))
	}
	if len(s.froms) > 0 {
		parts = append(parts, "FROM "+joinParts(s.froms, " ,"))
	}
	if len(s.joins) > 0 {
		parts = append(parts, joinParts(s.joins, " "))
	}
	if len(s.wheres) > 0 {
		parts = append(parts, "WHERE "+joinParts(s.wheres, " AND "))
	}
	if len(s.groups) > 0 {
		parts = append(parts, "GROUP BY "+joinParts(s.groups, " ,"))
	}
	if len(s.havings) > 0 {
		parts = append(parts, "HAVING "+joinParts(s.havings, " AND "))
	}
	if len(s.orders) > 0 {
		parts = append(parts, "ORDER BY "+joinParts(s.orders, " ,"))
	}
	if s.limit.partial != "" {
		parts = append(parts, "LIMIT "+s.limit.partial)
	}

	sql := strings.Join(parts, "\n")

	if s.PostgreSQL {
		sql = replacePositionalArguments(sql, 1)
	}
	return sql
}

func joinParts(components []component, joiner string) string {
	partials := make([]string, 0)
	for _, component := range components {
		partials = append(partials, component.partial)
	}
	return strings.Join(partials, joiner)
}

func appendArgs(args *[]interface{}, components ...component) {
	for _, component := range components {
		if len(component.args) > 0 {
			*args = append(*args, component.args...)
		}
	}
}

func replacePositionalArguments(sql string, c int) string {
	arg := fmt.Sprintf("$%d", c)
	newSql := strings.Replace(sql, "?", arg, 1)

	if newSql != sql {
		return replacePositionalArguments(newSql, c+1)
	}
	return newSql
}
