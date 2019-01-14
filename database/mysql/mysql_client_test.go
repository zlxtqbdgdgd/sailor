package mysql

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) { gc.TestingT(t) }

type mysqlSuite struct {
}

var _ = gc.Suite(&mysqlSuite{})

func (*mysqlSuite) TestFloat64(c *gc.C) {
	dsn := "root:@tcp(192.168.1.31:3306)/%s?timeout=1s&charset=utf8"
	dbt, err := GetClient(dsn, "gotest")
	if err != nil {
		c.Fatal(err)
	}
	dbt.Exec("DROP TABLE IF EXISTS test")

	types := [2]string{"FLOAT", "DOUBLE"}
	var expected float64 = 42.23
	var out float64
	var rows *sql.Rows
	for _, v := range types {
		dbt.Exec("CREATE TABLE test (value " + v + ")")
		dbt.Exec("INSERT INTO test VALUES (42.23)")
		rows, _ = dbt.Query("SELECT value FROM test")
		if rows.Next() {
			rows.Scan(&out)
			if expected != out {
				c.Fatalf("%s: %g != %g", v, expected, out)
			}
		} else {
			c.Fatalf("%s: no data", v)
		}
		dbt.Exec("DROP TABLE IF EXISTS test")
	}
}
