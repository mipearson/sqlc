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
	s = s.Join("LEFT JOIN Companies").Join("INNER JOIN Roles")

	expect(t, s.Args(), make([]interface{}, 0))
	expect(t, s.SQL(), strings.TrimSpace(`
SELECT *
FROM Employees
LEFT JOIN Companies INNER JOIN Roles
WHERE (name = 'Marge')
GROUP BY role
HAVING (a=1)
ORDER BY id
LIMIT 30
  `))
}

func TestArgumentComposition(t *testing.T) {
	s := Statement{}
	s = s.Where("name = ? OR name = ?", "Marge", "Alice").Where("role = ?", "Comptroller")
	expect(t, s.Args(), []interface{}{"Marge", "Alice", "Comptroller"})
	expect(t, s.SQL(), strings.TrimSpace("WHERE (name = ? OR name = ?) AND (role = ?)"))

	// PostgreSQL argument composition
	s.PostgreSQL = true
	expect(t, s.SQL(), strings.TrimSpace("WHERE (name = $1 OR name = $2) AND (role = $3)"))
}

func TestImmutability(t *testing.T) {
	orig := Statement{}
	orig = orig.Select("apples")
	expect(t, orig.SQL(), strings.TrimSpace("SELECT apples"))

	modified := orig.Select("oranges")
	expect(t, modified.SQL(), strings.TrimSpace("SELECT apples, oranges"))

	expect(t, orig.SQL(), strings.TrimSpace("SELECT apples"))
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("line %d: Got %#v, expected %#v", line, a, b)
	}
}
