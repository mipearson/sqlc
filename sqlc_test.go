package sqlc

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestBasicComposition(t *testing.T) {
	s := Statement{}
	// These statements are deliberately out of order
	s = s.Group("role").Order("id").Limit("30")
	s = s.Where("name = 'Marge'")
	s = s.Select("*").From("Employees")

	sql, args := s.ToSQL()
	expect(t, args, make([]interface{}, 0))
	expect(t, sql, strings.TrimSpace(`
SELECT *
FROM Employees
WHERE (name = 'Marge')
GROUP BY role
ORDER BY id
LIMIT 30
  `))
}

func TestArgumentComposition(t *testing.T) {
	s := Statement{}
	s = s.Where("name = ?", "Marge").Where("role = ?", "Comptroller")
	sql, args := s.ToSQL()
	expect(t, args, []interface{}{"Marge", "Comptroller"})
	expect(t, sql, strings.TrimSpace("WHERE (name = ?) AND (role = ?)"))

	// PostgreSQL argument composition
	s.PostgreSQL = true
	sql, _ = s.ToSQL()
	expect(t, sql, strings.TrimSpace("WHERE (name = $1) AND (role = $2)"))
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("line %d: Got %#v, expected %#v", line, a, b)
	}
}
