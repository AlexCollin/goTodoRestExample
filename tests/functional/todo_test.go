package functional

import (
	"bytes"
	"fmt"
	"github.com/AlexCollin/goTodoRestExample/model"
	"github.com/AlexCollin/goTodoRestExample/repo"
	"github.com/AlexCollin/goTodoRestExample/tests"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/AlexCollin/goTodoRestExample/handler"

	_ "github.com/lib/pq"
)

func TestGetAllTodo(t *testing.T) {
	repo := &repo.Todos{tests.Setup()}
	testServer := setupServer(repo)

	todo := &model.Todo{
		Title: "My Task1",
		Token: "123",
	}

	_, err := repo.Insert(todo)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8088/todo?token=123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	testServer.ServeHTTP(rec, req)

	got := strings.TrimSpace(rec.Body.String())

	want := `[{"id":1,"title":"My Task1","token":"123"}]`

	if len(got) == 0 {
		t.Fatalf("Want: %v, Got: %v", want, got)
	}
}

func TestInsertTodo(t *testing.T) {
	repo := &repo.Todos{tests.Setup()}
	testServer := setupServer(repo)

	body := []byte(`{"title":"My Task1","token":"123"}`)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8088/todo?token=123", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	testServer.ServeHTTP(rec, req)

	got := strings.TrimSpace(rec.Body.String())
	want := string(`{"id":1}`)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Want: %v, Got: %v", want, got)
	}

	gotTodo, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	wantTodo := []model.Todo{
		{
			Title: "My Task1",
			Token: "123",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Want: %v, Got: %v\n", wantTodo, gotTodo)
	}
}

func TestGetTodo(t *testing.T) {
	repo := &repo.Todos{tests.Setup()}
	testServer := setupServer(repo)

	todo := &model.Todo{
		Title: "My Task1",
		Token: "123",
	}

	id, err := repo.Insert(todo)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("http://localhost:8088/todo/%d?token=123", id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	testServer.ServeHTTP(rec, req)

	got := strings.TrimSpace(rec.Body.String())

	want := `{"id":1,"title":"My Task1","token":"123"}`

	if got != want {
		t.Fatalf("Want: %v, Got: %v", want, got)
	}
}

func TestUpdateTodo(t *testing.T) {
	repo := &repo.Todos{tests.Setup()}
	testServer := setupServer(repo)

	id, err := repo.Insert(&model.Todo{
		Title: "My Task1",
		Token: "123",
	})

	body := []byte(`{"title":"My Task2","token":"123"}`)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:8088/todo/%d?token=123", id), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	testServer.ServeHTTP(rec, req)

	got, err := repo.Get(id)
	if err != nil {
		t.Fatal(err)
	}

	want := model.Todo{
		ID:    1,
		Title: "My Task2",
		Token: "123",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Want: %v, Got: %v\n", got, want)
	}
}

func TestDeleteTodo(t *testing.T) {
	repo := &repo.Todos{tests.Setup()}
	testServer := setupServer(repo)

	id, err := repo.Insert(&model.Todo{
		Title: "My Task1",
		Token: "123",
	})

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8088/todo/%d?token=123", id), nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	testServer.ServeHTTP(rec, req)

	got, err := repo.Get(id)
	if got.ID > 0 {
		t.Fatalf("Should return the empty slice, Got: %v\n", got)
	}

}

func setupServer(repo *repo.Todos) *http.ServeMux {
	return handler.SetUpRouting(repo)
}
