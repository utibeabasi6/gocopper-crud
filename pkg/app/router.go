package app

import (
	"net/http"
	"strings"

	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/clogger"
	"github.com/utibeabasi6/gocopper-crud/pkg/todos"
)

type NewRouterParams struct {
	Todos  *todos.Queries
	RW     *chttp.ReaderWriter
	Logger clogger.Logger
}

func NewRouter(p NewRouterParams) *Router {
	return &Router{
		rw:     p.RW,
		logger: p.Logger,
		todos:  p.Todos,
	}
}

type Router struct {
	rw     *chttp.ReaderWriter
	logger clogger.Logger
	todos  *todos.Queries
}

func (ro *Router) Routes() []chttp.Route {
	return []chttp.Route{
		{
			Path:    "/",
			Methods: []string{http.MethodGet},
			Handler: ro.HandleIndexPage,
		},
		{
			Path:    "/todos",
			Methods: []string{http.MethodGet},
			Handler: ro.HandleTodosPage,
		},
		{
			Path:    "/todos",
			Methods: []string{http.MethodPost},
			Handler: ro.HandleCreateTodos,
		},
		{
			Path:    "/{todo}",
			Methods: []string{http.MethodPost},
			Handler: ro.HandleUpdateTodos,
		},
		{
			Path:    "/{todo}",
			Methods: []string{http.MethodDelete},
			Handler: ro.HandleDeleteTodos,
		},
	}
}

func (ro *Router) HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	ro.rw.WriteHTML(w, r, chttp.WriteHTMLParams{
		PageTemplate: "index.html",
	})
}

func (ro *Router) HandleTodosPage(w http.ResponseWriter, r *http.Request) {
	todos, err := ro.todos.ListTodos(r.Context())
	if err != nil {
		ro.logger.Error("an error occured while fetching todos", err)
		ro.rw.WriteHTMLError(w, r, cerrors.New(nil, "unable to fetch todos", map[string]interface{}{
			"form": r.Form,
		}))
	}
	ro.rw.WriteHTML(w, r, chttp.WriteHTMLParams{
		PageTemplate: "todos.html",
		Data:         todos,
	})
}

func (ro *Router) HandleCreateTodos(w http.ResponseWriter, r *http.Request) {
	var (
		todo = strings.TrimSpace(r.PostFormValue(("todo")))
	)
	if todo == "" {
		ro.rw.WriteHTMLError(w, r, cerrors.New(nil, "unable to create todos: todo cannot be empty", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	newtodo := todos.Todo{Name: todo}
	err := ro.todos.SaveTodo(r.Context(), &newtodo)
	if err != nil {
		ro.logger.Error("an error occured while saving todo", err)
		ro.rw.WriteHTMLError(w, r, cerrors.New(nil, "unable to create todos", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}
	http.Redirect(w, r, "/todos", http.StatusSeeOther)
}

type error struct {
	error string
}

func (ro *Router) HandleDeleteTodos(w http.ResponseWriter, r *http.Request) {
	var (
		todo = strings.TrimSpace(chttp.URLParams(r)["todo"])
	)
	if todo == "" {
		deleteError := error{error: "Unable to delete todo"}
		ro.rw.WriteJSON(w, chttp.WriteJSONParams{StatusCode: 500, Data: deleteError})
		return
	}

	newtodo := todos.Todo{Name: todo}
	err := ro.todos.DeleteTodo(r.Context(), &newtodo)
	if err != nil {
		ro.logger.Error("an error occured while deleting todo", err)
		deleteError := error{error: "Unable to delete todo"}
		ro.rw.WriteJSON(w, chttp.WriteJSONParams{StatusCode: 500, Data: deleteError})
		return
	}
	http.Redirect(w, r, "/todos", http.StatusSeeOther)
}

func (ro *Router) HandleUpdateTodos(w http.ResponseWriter, r *http.Request) {
	var (
		todo    = strings.TrimSpace(r.PostFormValue("todo"))
		oldName = chttp.URLParams(r)["todo"]
	)
	if todo == "" {
		ro.rw.WriteHTMLError(w, r, cerrors.New(nil, "unable to update todos: todo cannot be empty", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}

	newtodo := todos.Todo{Name: todo}
	ro.logger.WithTags(map[string]interface{}{
		"oldname": oldName,
		"newname": newtodo.Name,
	}).Info("Todo updated")
	err := ro.todos.UpdateTodo(r.Context(), oldName, &newtodo)
	if err != nil {
		ro.logger.Error("an error occured while saving todo", err)
		ro.rw.WriteHTMLError(w, r, cerrors.New(nil, "unable to update todos", map[string]interface{}{
			"form": r.Form,
		}))
		return
	}
	http.Redirect(w, r, "/todos", http.StatusSeeOther)
}
