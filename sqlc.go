// Package sqlc helps you combine bits of arbitrary bits of SQL
package sqlc

import "strings"

type component struct {
	partial string
	args    []interface{}
}

// From the psql documentation:
// SELECT [ ALL | DISTINCT [ ON ( expression [, ...] ) ] ]
//     * | expression [ [ AS ] output_name ] [, ...]
//     [ FROM from_item [, ...] ]
//     [ WHERE condition ]
//     [ GROUP BY expression [, ...] ]
//     [ HAVING condition [, ...] ]
//     [ WINDOW window_name AS ( window_definition ) [, ...] ]
//     [ { UNION | INTERSECT | EXCEPT } [ ALL | DISTINCT ] select ]
//     [ ORDER BY expression [ ASC | DESC | USING operator ] [ NULLS { FIRST | LAST } ] [, ...] ]
//     [ LIMIT { count | ALL } ]
//     [ OFFSET start [ ROW | ROWS ] ]
//     [ FETCH { FIRST | NEXT } [ count ] { ROW | ROWS } ONLY ]
//     [ FOR { UPDATE | NO KEY UPDATE | SHARE | KEY SHARE } [ OF table_name [, ...] ] [ NOWAIT ] [...] ]

type Statement struct {
	selects []component
	froms   []component
	wheres  []component
	groups  []component
	orders  []component
	limit   component
}

func (s Statement) Select(partial string, args ...interface{}) Statement {
	s.selects = append(s.selects, component{partial, args})
	return s
}

func (s Statement) From(partial string, args ...interface{}) Statement {
	s.froms = append(s.froms, component{partial, args})
	return s
}

func (s Statement) Where(partial string, args ...interface{}) Statement {
	s.wheres = append(s.wheres, component{"(" + partial + ")", args})
	return s
}

func (s Statement) Group(partial string, args ...interface{}) Statement {
	s.groups = append(s.groups, component{partial, args})
	return s
}

func (s Statement) Order(partial string, args ...interface{}) Statement {
	s.orders = append(s.orders, component{partial, args})
	return s
}

func (s Statement) Limit(partial string, args ...interface{}) Statement {
	s.limit = component{partial, args}
	return s
}

func (s Statement) ToSQL() (string, []interface{}) {
	parts := make([]string, 0)
	args := make([]interface{}, 0)

	if len(s.selects) > 0 {
		parts = append(parts, "SELECT "+joinParts(s.selects, ", ", &args))
	}
	if len(s.froms) > 0 {
		parts = append(parts, "FROM "+joinParts(s.froms, " ,", &args))
	}
	if len(s.wheres) > 0 {
		parts = append(parts, "WHERE "+joinParts(s.wheres, " AND ", &args))
	}
	if len(s.groups) > 0 {
		parts = append(parts, "GROUP BY "+joinParts(s.groups, " ,", &args))
	}
	if len(s.orders) > 0 {
		parts = append(parts, "ORDER BY "+joinParts(s.orders, " ,", &args))
	}
	if s.limit.partial != "" {
		parts = append(parts, "LIMIT "+addPart(s.limit, &args))
	}

	return strings.Join(parts, "\n"), args
}

func joinParts(components []component, joiner string, args *[]interface{}) string {
	partials := make([]string, 0)
	for _, component := range components {
		partials = append(partials, addPart(component, args))
	}
	return strings.Join(partials, joiner)
}

func addPart(component component, args *[]interface{}) string {
	*args = append(*args, component.args...)
	return component.partial
}
