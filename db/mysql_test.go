package db

import (
	"fmt"

	"testing"
)

func TestDB(t *testing.T) {
	rows, err := Conn().Query("select * from sys_user")
	fmt.Println(rows.Columns())
	fmt.Println(err)
}
