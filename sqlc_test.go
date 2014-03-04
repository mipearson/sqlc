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
	s = s.Having("a=1")

	sql, args := s.ToSQL()
	expect(t, args, make([]interface{}, 0))
	expect(t, sql, strings.TrimSpace(`
SELECT *
FROM Employees
WHERE (name = 'Marge')
GROUP BY role
HAVING (a=1)
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

func TestImmutability(t *testing.T) {
	orig := Statement{}
	orig = orig.Select("apples")
	sql, _ := orig.ToSQL()
	expect(t, sql, strings.TrimSpace("SELECT apples"))

	modified := orig.Select("oranges")
	sql, _ = modified.ToSQL()
	expect(t, sql, strings.TrimSpace("SELECT apples, oranges"))

	sql, _ = orig.ToSQL()
	expect(t, sql, strings.TrimSpace("SELECT apples"))
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("line %d: Got %#v, expected %#v", line, a, b)
	}
}
