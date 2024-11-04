package main

import(
	"fmt"
	"html/template"
	"os"
	"net/http"
	"log"
	"regexp"
	//"errors"
)

type Page struct {
	Title string
	Body [] byte
}

func (p * Page) save() error{
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body,0600)
}

func loadPage(title string)(* Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body:body}, nil
}

func handler ( w http.ResponseWriter, r * http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func viewHandler (w http.ResponseWriter, r *http.Request, title string){
	//title,err := getTitle(w,r)
	
	//if err!= nil {
	//	return
	//}
	p,err := loadPage(title)
	if err != nil {
		http.Redirect(w,r, "/edit/"+ title, http.StatusFound)
		return
	}
	//fmt.Fprintf(w, "<h1>%s<h1><div>%s</div>", p.Title, p.Body)
	renderTemplate(w, "view",p)
}


func editHandler(w http.ResponseWriter, r *http.Request, title string) {

	//title,err := getTitle(w,r)
	
	//if err!= nil {
	//	return
	//}
	p, err := loadPage(title)

	if err != nil {

		p = &Page{Title: title}

	}

	renderTemplate(w, "edit", p)
	//fmt.FPrintf(w, "<h1> Editing %s </h1>" + 
	//"<form action=\"/save/%s\" method=\"POST\">" +
	//"<textarea name=\"body\"> %s </textarea><br>" +
	//"<input type=\"submit\" value=\"Save\">" +
	//"</form>",
	//p.Title, p.Title, p.Body)

}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	
	// title, err := getTitle(w, r)
	// if err != nil {
	//     return
	// }

	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}

	// Declare err here to be used for the save function's return value
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseFiles("edit.html","view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p* Page) {
	
	err := templates.ExecuteTemplate(w, tmpl+".html",p)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
}


var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
func makeHandler (fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m:=validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w,r)
			return
		}
		fn(w,r,m[2])
	}
}
func main () {

	//http.HandleFunc("/", handler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080",nil))
	//p1 := &Page{Title : "TestPage", Body: []byte("This is a simple Page.")}
	//p1.save()
	//p2, _ := loadPage("TestPage")
	//fmt.Println(string(p2.Body))
}