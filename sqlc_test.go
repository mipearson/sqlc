package sqlc

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestBasicComposition(t *testing.T) {
	s := Statement{}
	s = s.Select("*").From("Employees").Where("name = 'Marge'").Order("id")

	sql, args := s.ToSQL()
	expect(t, args, make([]interface{}, 0))
	expect(t, sql, strings.TrimSpace(`
SELECT *
FROM Employees
WHERE (name = 'Marge')
ORDER BY id
  `))
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("line %d: Got %#v, expected %#v", line, a, b)
	}
}
