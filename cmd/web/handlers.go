package main

import (
	"Todo_Application/pkg/forms"
	"Todo_Application/pkg/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

var Task_array []*models.Todo
var Forms *forms.Form
var errors = make(map[string]string)

func (a *application) App(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { //path not found
		a.notFound(w)
		return
	}
	s, errr := a.Todo.GetAll() //Getting all the data from db
	if errr != nil {
		log.Println(errr)
		a.serverError(w, errr)
		http.Error(w, "Internal Server Error1", 500)
		return
	}
	//Flash := a.session.PopString(r, "flash")
	//panic("oops!something went wrong")
	tmpl, err := template.ParseFiles("./ui/html/App.page.tmpl") //parsing the template file
	if err != nil {
		log.Println(err)
		a.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error1", 500)
		return
	}
	//flash := a.session.GetString(r, "flash")
	err = tmpl.Execute(w, struct {
		Tasks []*models.Todo
		Flash string
	}{
		Tasks: s,
		Flash: a.session.PopString(r, "flash"),
	}) //executing the template
	if err != nil {
		log.Println((err.Error()))
		http.Error(w, "Internal Server Error", 500)
	} else {
		a.infoLog.Println("App page Reloaded")
	}
	a.session.Put(r, "flash", "Todo successfully created!")
}

func (a *application) Add(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("todoText")

	//errors := make(map[string]string)

	if strings.TrimSpace(name) == "" { //checking user input string
		errors["name"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(name) > 100 {
		errors["name"] = "This field is too long (maximum is 100 characters)"
	}

	if len(errors) > 0 {
		a.session.Put(r, "flash", "This field cannot be blank")
		s, err := a.Todo.ErrorManage(errors)
		log.Println(s)
		if err != nil {
			a.errorLog.Println(err.Error())
			return
		}
		tmpl, err := template.ParseFiles("./ui/html/App.page.tmpl")
		if err != nil {
			a.errorLog.Println(err.Error())
			return
		}
		Task_array = append(Task_array, s)
		//fmt.Print(errors)
		err = tmpl.Execute(w, Task_array)
		if err != nil {
			log.Println((err.Error()))
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}
	a.session.Put(r, "flash", "The add session is created")
	_, err := a.Todo.Insert(name, "365") //inserting value from form to insert method of todo
	if err != nil {
		a.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *application) Update(w http.ResponseWriter, r *http.Request) {
	id1, _ := strconv.Atoi(r.FormValue("id"))
	//name1 := r.FormValue("update_name")
	_, err := a.Todo.Update(r.FormValue("update_name"), id1) //calling update function
	if err != nil {
		a.errorLog.Println(err)
		return
	}
	a.session.Put(r, "flash", "The session is updated")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *application) Delete(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(r.FormValue("ID")) //Getting the id from form to delete
	_, err := a.Todo.Delete(id)              //Calling the delete function
	if err != nil {
		a.errorLog.Println(err)
		return
	}
	a.session.Put(r, "flash", "The delete session is running")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *application) signupUserForm(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("./ui/html/signup.page.tmpl")
	if err != nil {
		a.errorLog.Println(err.Error())
		return
	}
	f := Forms.New(nil)
	err = tmpl.Execute(w, f)
	if err != nil {
		log.Println((err.Error()))
		http.Error(w, "Internal Server Error2", 500)

	}
}

func (a *application) signupUser(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Create a new user...")
	err := r.ParseForm()
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}
	form := Forms.New(r.PostForm)
	form.Required("names", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	errr := a.users.Insert(form.Get("names"), form.Get("email"), form.Get("password"))
	if errr != nil {
		log.Println((errr.Error()))
		http.Error(w, "Internal Server Error2", 500)

	}
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		// err = tmpl.Execute(w, form)
		// if err != nil {
		// 	log.Println((err.Error()))
		// 	http.Error(w, "Internal Server Error2", 500)
		// }
		a.session.Put(r, "flash", "Your signup was successful. Please log in.")
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}

func (a *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "Display the user login form...")
	tmpl, err := template.ParseFiles("./ui/html/login.page.tmpl")
	if err != nil {
		a.errorLog.Println(err.Error())
		return
	}
	f := Forms.New(nil)
	er := tmpl.Execute(w, f)
	if er != nil {
		log.Println((er.Error()))
		http.Error(w, "Internal Server Error3", 500)
	}
}

func (a *application) loginUser(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "Authenticate and login the user...")
	err := r.ParseForm()
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}
	form := Forms.New(r.PostForm)
	id, err := a.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		a.session.Put(r, "flash", "Email or Password is incorrect! ")
	} else if err != nil {
		a.serverError(w, err)
		return
	}
	a.session.Put(r, "userID", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userID")
	app.session.Put(r, "flash", "You've been logged out successfully! ")
	http.Redirect(w, r, "/", 303)
}
