# Go Web App Task - Pemrograman Berbasis Kerangka Kerja

## Nama : Helsa Sriprameswari Putri
## NRP : 5025221154
## Kelas : D

## Page Struct

```
type Page struct {

	Title string
	Body  []byte

}
```

Membuat struct bernama Page untuk menyimpan title berupa string dan body berupa slice yang berisi kata - kata.

## Function Save
```
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

```

Fungsi untuk menyimpan isi body ke dalam title (berupa file txt) dengan read and write only permission.

## Function loadPage
```
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}
```

Fungsi loadPage digunakan untuk menampilkan suatu page berisikan title dan body dengan membaca file txt yang sudah exist.

## Function viewHandler

```
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", p)
}
```

Fungsi ini digunakan untuk menghandle page yang mau di view. Jika title tidak exist, page akan meminta user untuk mengedit title baru tersebut dan membuatnya. Jika title/page exist, maka akan mereturn view sesuai page yang diminta

## Function editHandler
```
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)

	if err != nil {
		p = &Page{Title: title}

	}
	renderTemplate(w, "edit", p)
}
```

Fungsi editHandler digunakan untuk mengedit page dengan isi body baru dengan loadPage tersebut lalu baru bisa diedit.

## Function saveHandler
```
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {

	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
```
Fungsi saveHandler untuk menyimpan page dengan menerima formValue yang berisi body/isi page nya. Jika title dan body page sukses disimpan, akan mereturn view dengan title tersebut.

## Function renderTemplates

```
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
```
Membuat variable templates untuk mengakses file edit.html dan view.html yang sudah dibuat. Fungsi renderTemplate ini untuk mengeksekusi template html sesuai yang diminta.

## Function makeHandler
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
```
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)

		if m == nil {
			http.NotFound(w, r)
			return
		}

		fn(w, r, m[2])

	}

}
```

Membuat variable validPath untuk menerima path tertentu yang sesuai harus ada dalam /edit/ atau /save/ atau /view/. Jika user menginputkan path lain, akan return error page not found. Jika path sesuai, akan memproses page sesuai titlenya.

## Function main
```
func main() {

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Fungsi main untuk memproses dan menampilkan halaman serta mengatur routing view, edit, save  dalam server http port 8080


## Web App Overview

![image](https://github.com/user-attachments/assets/332320de-bc63-4486-a629-11f69ec94fca)

Saat mengakses http://localhost:8080/view/test, web akan menampilkan title "test" dengan b ody/isi sesuai yang ditulis di test.txt.

![image](https://github.com/user-attachments/assets/a8530516-4ce7-4376-b1d5-37da37dfbd21)

Halaman ini untuk mengedit body/isi dari test page dalam path /edit/test.

![image](https://github.com/user-attachments/assets/b7df487e-2bd6-4dee-8440-ed2662b67e21)

Jika halaman yang di edit sudah disave, akan return ke view page tersebut dengan isi/body baru sesuai yang sudah diedit.

![image](https://github.com/user-attachments/assets/98285e9c-7bcc-4bc4-8924-33af22f77631)

Halamann TestPage juga bisa diview dan menampilkan isi sesuai TestPage.txt. Bisa juga untuk diedit dan disave seperti halaman test.

![image](https://github.com/user-attachments/assets/48a07671-9fcf-48a2-b9a5-f245a5f5c4a2)

Jika mencoba untuk view halaman baru (belum exist), contohnya ke http://localhost:8080/view/hello, halaman akan return ke view http://localhost:8080/edit/hello untuk menambahkan halaman hello dengan membuat body pagenya.

![image](https://github.com/user-attachments/assets/7bfa22fd-fde6-4c29-b3e4-c5389c9aca5f)

Setelah halaman baru diedit dan save, akan terbentuk halaman baru bernama hello sesuai body yang kita tulis dan sudah dapat di view halamannya. File hello.txt juga muncul dalam direktori local.

![image](https://github.com/user-attachments/assets/d9b9a6fd-68e6-49a7-9e0e-340cccd818b9)

Tampilan web apabila mencoba akses dengan path yang tidak sesuai, akan mereturn 404 page not found.
