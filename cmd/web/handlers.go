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
	tmpl, err := template.ParseFiles("./ui/html/App.page.tmpl") //parsing the template file
	if err != nil {
		log.Println(err)
		a.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error1", 500)
		return
	}
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
	if strings.TrimSpace(name) == "" { //checking user input string
		errors["name"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(name) > 100 {
		errors["name"] = "This field is too long (maximum is 100 characters)"
	}
	//Input validation module
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
		err = tmpl.Execute(w, Task_array)
		if err != nil {
			log.Println((err.Error()))
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}
	//Input validation ends
	a.session.Put(r, "flash", "The add session is created")
	ok := strings.Contains(name, "special")
	types := "normal"
	if ok {
		types = "special"
		_, err := a.Special.Insert(name, types, "365") //inserting value from form to insert method of special
		if err != nil {
			a.serverError(w, err)
			return
		}
	}
	_, err := a.Todo.Insert(name, types, "365") //inserting value from form to insert method of todo
	if err != nil {
		a.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *application) Update(w http.ResponseWriter, r *http.Request) {
	id1, _ := strconv.Atoi(r.FormValue("id"))
	_, err := a.Todo.Update(r.FormValue("update_name"), id1) //calling update function
	if err != nil {
		a.errorLog.Println(err)
		return
	}
	a.session.Put(r, "flash", "The session is updated")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *application) Delete(w http.ResponseWriter, r *http.Request) {
	names := r.FormValue("Name")   //Getting the name from form to delete
	_, err := a.Todo.Delete(names) //Calling the delete function
	if err != nil {
		a.errorLog.Println(err)
		return
	}
	_, errr := a.Special.Delete(names) //Calling the delete function
	if errr != nil {
		return
	}
	a.session.Put(r, "flash", "The delete session is running")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *application) specialDelete(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("Name")      //Getting the name from form to delete
	_, err := a.Special.Delete(name) //Calling the delete function
	if err != nil {
		a.errorLog.Println(err)
		return
	}
	_, errr := a.Todo.Delete(name) //Calling the delete function
	if errr != nil {
		a.errorLog.Println(errr)
		return
	}
	http.Redirect(w, r, "/special", http.StatusSeeOther)
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
		a.session.Put(r, "flash", "Your signup was successful. Please log in.")
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}

func (a *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
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

func (a *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	a.session.Remove(r, "userID")
	a.session.Put(r, "flash", "You've been logged out successfully! ")
	http.Redirect(w, r, "/", 303)
}

func (a *application) special(w http.ResponseWriter, r *http.Request) {
	s, err := a.Special.GetSpecial() //Getting the special task data from db
	if err != nil {
		//log.Println(err)
		a.serverError(w, err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	tmpl, err := template.ParseFiles("./ui/html/special.page.tmpl") //parsing the template file
	if err != nil {
		//log.Println(err)
		a.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error1", 500)
		return
	}
	err = tmpl.Execute(w, struct{ Task_specials []*models.Special }{Task_specials: s}) //executing the template
	if err != nil {
		log.Println((err.Error()))
		http.Error(w, "Internal Server Error", 500)
	} else {
		a.infoLog.Println("Special page Running")
	}
}
