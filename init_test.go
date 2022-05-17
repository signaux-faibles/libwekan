package libwekan

import (
	"context"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	w, err := Connect(context.Background(), "localhost:27017", "test", "test")
	if err != nil {
		t.Fatal("testing env not suitable for test: \nlibwekan => " + err.Error())
	}

	fmt.Println(w.ListUsers(context.Background()))
}
