package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"html/template"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"

	"apitest.com/api/controllers"
	"apitest.com/api/models"
)

var (
	BlogBaseRegex          = regexp.MustCompile(`^/blogs/*$`)
	BlogByIdRegex          = regexp.MustCompile(`^/blogs/([0-9]+)`)
	BlogByNameRegex        = regexp.MustCompile(`^/blogs/name/([a-zA-Z0-9_.-]*)$`)
	BlogFileByIdRegex      = regexp.MustCompile(`^/blogs/file/([0-9]+)`)
	BlogForm               = regexp.MustCompile(`^/blogs/admin/fileform/*$`)
	AdminBlogBaseRegex     = regexp.MustCompile(`^/blogs/admin/*$`)
	AdminBlogByIdRegex     = regexp.MustCompile(`^/blogs/admin/([0-9]+)`)
	AdminBlogByNameRegex   = regexp.MustCompile(`^/blogs/admin/name/([a-zA-Z0-9_.-]*)$`)
	AdminBlogFileByIdRegex = regexp.MustCompile(`^/blogs/admin/file/([0-9]+)`)
)

type BlogsRoute struct{}

// main router function
func (d *BlogsRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	// ADMIN
	case r.Method == http.MethodPost && AdminBlogBaseRegex.MatchString(r.URL.Path):
		//d.AdminCreateUser(w, r)
		return
	case r.Method == http.MethodGet && AdminBlogBaseRegex.MatchString(r.URL.Path):
		//d.AdminGetUsers(w, r)
		return
	case r.Method == http.MethodGet && AdminBlogByIdRegex.MatchString(r.URL.Path):
		// put admin method here
		return
	case r.Method == http.MethodGet && AdminBlogByNameRegex.MatchString(r.URL.Path):
		// put admin method hered
		return
	case r.Method == http.MethodPut && AdminBlogByIdRegex.MatchString(r.URL.Path):
		// put admin method here
		return
	case r.Method == http.MethodDelete && AdminBlogByIdRegex.MatchString(r.URL.Path):
		// put admin method here
		return

		// NORMAL
	case r.Method == http.MethodPost && BlogBaseRegex.MatchString(r.URL.Path):
		d.CreateBlog(w, r)
		return
	case r.Method == http.MethodGet && BlogBaseRegex.MatchString(r.URL.Path):
		d.GetBlogs(w, r)
		return
	case r.Method == http.MethodGet && BlogByIdRegex.MatchString(r.URL.Path):
		d.GetBlogById(w, r)
		return
	case r.Method == http.MethodGet && BlogByNameRegex.MatchString(r.URL.Path):
		d.GetBlogByName(w, r)
		return
	case r.Method == http.MethodPut && BlogByIdRegex.MatchString(r.URL.Path):
		//d.UpdateUserKeyValuePair(w, r)
		return
	case r.Method == http.MethodDelete && BlogByIdRegex.MatchString(r.URL.Path):
		d.DeleteBlogById(w, r)
		return
	case r.Method == http.MethodPost && BlogFileByIdRegex.MatchString(r.URL.Path):
		d.receiveMDFile(w, r)
		return
	case r.Method == http.MethodGet && BlogFileByIdRegex.MatchString(r.URL.Path):
		d.serveFile(w, r)
		return
    case r.Method == http.MethodPost && BlogForm.MatchString(r.URL.Path):
        d.sendForm(w, r)
		return

	default:
		w.Write([]byte("Test response from tags"))
		return
	}
}

// FRONTEND ROUTES
func (d *BlogsRoute) CreateBlog(w http.ResponseWriter, r *http.Request) {
	var b models.Blog
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while decoding blog post from json", err)
		return
	}

	newBlog, err := b.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while creating blog: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	jsonBlog, err := json.Marshal(newBlog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while marshalling blog to json: ", err)
		return
	}

	w.Write(jsonBlog)
}

// this route takes in id from link
func (d *BlogsRoute) receiveMDFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error during file upload: ", err)
		return
	}
	defer file.Close()

    root, err := os.Getwd()
    mdpath := filepath.Join(root, "views", "mdfiles", header.Filename)

	endFile, err := os.Create(mdpath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error during file save to markdown directory: ", err)
		return
	}
	defer endFile.Close()

	var data bytes.Buffer
    io.Copy(&data, file)
    contents := data.String()

	_, err = endFile.WriteString(contents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error during data writing into markdown file: ", err)
		return
	}

	htmlpath, err := controllers.MdToHTML(mdpath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error during data writing into markdown file: ", err)
		return
	}

	id_path := BlogFileByIdRegex.FindStringSubmatch(r.URL.Path)
	bid, err := strconv.Atoi(id_path[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while getting blog id from path", err)
		return
	}

	blog, err := models.ReadBlogById(bid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while getting blog from the bid given", err)
		return
	}

	err = blog.UpdateLinks(mdpath, htmlpath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while updating blog links", err)
		return
	}

    w.Header().Set("Content-Type", "text/html")

	templateFile := htmlpath
	tmpl, err := template.New(filepath.Base(templateFile)).ParseFiles(templateFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while parsing template in blog", err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while executing template in blog", err)
		return
	}
}

func (d *BlogsRoute) serveFile(w http.ResponseWriter, r *http.Request) {
	id_path := BlogFileByIdRegex.FindStringSubmatch(r.URL.Path)
	bid, err := strconv.Atoi(id_path[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while getting blog id from path", err)
		return
	}

	blog, err := models.ReadBlogById(bid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while getting blog from the bid given", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	templateFile := blog.Link_to_jsx
	tmpl, err := template.New(filepath.Base(templateFile)).ParseFiles(templateFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while parsing template in blog", err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while executing template in blog", err)
		return
	}
}

func (d *BlogsRoute) sendForm(w http.ResponseWriter, r *http.Request){
    r.ParseForm()

    bid := r.FormValue("blog-id")

	w.Header().Set("Content-Type", "text/html")

	templateFile := "views/blog_file_form.html"
	tmpl, err := template.New(filepath.Base(templateFile)).ParseFiles(templateFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while parsing template in blog", err)
		return
	}

    type B struct{
        Id int
    }
    var b B
    b.Id, err = strconv.Atoi(bid)
    if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error conveting bid to integer in blog form sending", err)
		return
	}

	err = tmpl.Execute(w, b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while executing template in blog", err)
		return
	}
}

func (d *BlogsRoute) GetBlogs(w http.ResponseWriter, r *http.Request) {
	blogSlice, err := models.ReadAllBlogs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while getting blog list", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(blogSlice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while marshalling blog slice to json: ", err)
		return
	}

	w.Write(jsonData)
}

func (d *BlogsRoute) GetBlogById(w http.ResponseWriter, r *http.Request) {
	id_path := BlogByIdRegex.FindStringSubmatch(r.URL.Path)
	bid, err := strconv.Atoi(id_path[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while getting id from path", err)
		return
	}

	b, err := models.ReadBlogById(bid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while reading blog by id", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while marshalling blog to json while getting blog by id: ", err)
		return
	}

	w.Write(jsonData)
}

func (d *BlogsRoute) GetBlogByName(w http.ResponseWriter, r *http.Request) {
	id_path := BlogByNameRegex.FindStringSubmatch(r.URL.Path)
	blogTitle := id_path[1]

	b, err := models.ReadTagByName(blogTitle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while reading blog by name", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while marshalling blog to json while getting blog by name: ", err)
		return
	}

	w.Write(jsonData)
}

func (d *BlogsRoute) DeleteBlogById(w http.ResponseWriter, r *http.Request) {
	id_path := BlogByIdRegex.FindStringSubmatch(r.URL.Path)
	bid, err := strconv.Atoi(id_path[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while getting id from path during blog deletion", err)
		return
	}

	b, err := models.ReadBlogById(bid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while reading blog by id during deletion", err)
		return
	}

	err = b.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while deleting blog", err)
		return
	}

	w.Write([]byte(fmt.Sprintf("Blog with id %d was deleted succesfully", bid)))
}
