package user

import (
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/asdine/storm"
	"gopkg.in/mgo.v2/bson"
)

func TestMain(m *testing.M) {
	m.Run()
	os.Remove(dbPath)
}

func cleanDB(b *testing.B) {
	os.Remove(dbPath)
	u := &User{
		ID:   bson.NewObjectId(),
		Name: "Jhon",
		Role: "Tester",
	}
	err := u.Save()
	if err != nil {
		b.Fatalf("Error saving a record: %s", err)
	}
	b.ResetTimer()
}

func BenchmarkCreate(b *testing.B) {
	cleanDB(b)
	os.Remove(dbPath)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Jhon_" + strconv.Itoa(i),
			Role: "Tester",
		}
		b.StartTimer()
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
	}
}

func BenchmarkRead(b *testing.B) {
	cleanDB(b)
	os.Remove(dbPath)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Jhon_" + strconv.Itoa(i),
			Role: "Tester",
		}
		b.StartTimer()
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
		_, err = One(u.ID)
		if err != nil {
			b.Fatalf("Error retriving record %s", err)
		}
	}
}

func BenchmarkUpdate(b *testing.B) {
	cleanDB(b)
	os.Remove(dbPath)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Jhon_" + strconv.Itoa(i),
			Role: "Tester",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
		u.Role = "Developer"
		b.StartTimer()
		err = u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
		_, err = One(u.ID)
		if err != nil {
			b.Fatalf("Error retriving record %s", err)
		}
	}
}

func BenchmarkDelete(b *testing.B) {
	cleanDB(b)
	os.Remove(dbPath)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Jhon_" + strconv.Itoa(i),
			Role: "Tester",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
		b.StartTimer()
		err = Delete(u.ID)
		if err != nil {
			b.Fatalf("Error removing record %s", err)
		}
		_, err = One(u.ID)
		if err == nil {
			b.Fatal("Record should not exist anymore")
		}
	}
}

func BenchmarkCRUD(b *testing.B) {
	os.Remove(dbPath)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Jhon",
			Role: "Tester",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}

		_, err = One(u.ID)
		if err != nil {
			b.Fatalf("Error retriving record %s", err)
		}

		u.Role = "Developer"
		err = u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
		_, err = One(u.ID)
		if err != nil {
			b.Fatalf("Error retriving record %s", err)
		}

		err = Delete(u.ID)
		if err != nil {
			b.Fatalf("Error removing record %s", err)
		}
		_, err = One(u.ID)
		if err == nil {
			b.Fatal("Record should not exist anymore")
		}
		if err != storm.ErrNotFound {
			b.Fatalf("Error retriving non-existing record %s", err)
		}
	}

}
func TestCRUD(t *testing.T) {
	t.Log("Create")
	u := &User{
		ID:   bson.NewObjectId(),
		Name: "Jhon",
		Role: "Tester",
	}
	err := u.Save()
	if err != nil {
		t.Fatalf("Error saving a record: %s", err)
	}

	t.Log("Read")
	u2, err := One(u.ID)
	if err != nil {
		t.Fatalf("Error retriving record %s", err)
	}
	if !reflect.DeepEqual(u2, u) {
		t.Error("Records do not match")
	}

	t.Log("Update")
	u.Role = "Developer"
	err = u.Save()
	if err != nil {
		t.Fatalf("Error saving a record: %s", err)
	}
	u3, err := One(u.ID)
	if err != nil {
		t.Fatalf("Error retriving record %s", err)
	}
	if !reflect.DeepEqual(u3, u) {
		t.Error("Records do not match")
	}
	t.Log("Delete")
	err = Delete(u.ID)
	if err != nil {
		t.Fatalf("Error removing record %s", err)
	}
	_, err = One(u.ID)
	if err == nil {
		t.Fatal("Record should not exist anymore")
	}
	if err != storm.ErrNotFound {
		t.Fatalf("Error retriving non-existing record %s", err)
	}

	t.Log("Read All")
	u2.ID = bson.NewObjectId()
	u3.ID = bson.NewObjectId()
	err = u2.Save()
	if err != nil {
		t.Fatalf("Error saving a record: %s", err)
	}
	err = u3.Save()
	if err != nil {
		t.Fatalf("Error saving a record: %s", err)
	}
	users, err := All()
	if err != nil {
		t.Fatalf("Error reading all records %s", err)
	}
	if len(users) != 2 {
		t.Errorf("Different number of records retrieved. %d", len(users))
	}
}
