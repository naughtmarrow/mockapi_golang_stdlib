package routes

import (
	"encoding/json"
	"fmt"
	//"html/template"
	"net/http"
	//"path/filepath"
	"regexp"
	"strconv"

	"apitest.com/api/models"
)

var (
	TagBaseRegex        = regexp.MustCompile(`^/tags/*$`)
	TagByIdRegex        = regexp.MustCompile(`^/tags/([0-9]+)`)
	TagByNameRegex      = regexp.MustCompile(`^/tags/name/([a-zA-Z0-9_.-]*)$`)
	AdminTagBaseRegex   = regexp.MustCompile(`^/tags/admin/*$`)
	AdminTagByIdRegex   = regexp.MustCompile(`^/tags/admin/([0-9]+)`)
	AdminTagByNameRegex = regexp.MustCompile(`^/tags/admin/name/([a-zA-Z0-9_.-]*)$`)
)

type TagsRoute struct{}

// main router function
func (d *TagsRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
    // ADMIN
    case r.Method == http.MethodPost && AdminTagBaseRegex.MatchString(r.URL.Path):
        //d.AdminCreateUser(w, r)
		return
	case r.Method == http.MethodGet && AdminTagBaseRegex.MatchString(r.URL.Path):
		//d.AdminGetUsers(w, r)
		return
	case r.Method == http.MethodGet && AdminTagByIdRegex.MatchString(r.URL.Path):
		// put admin method here
		return
	case r.Method == http.MethodGet && AdminTagByNameRegex.MatchString(r.URL.Path):
		// put admin method hered
		return
	case r.Method == http.MethodPut && AdminTagByIdRegex.MatchString(r.URL.Path):
		// put admin method here
		return
	case r.Method == http.MethodDelete && AdminTagByIdRegex.MatchString(r.URL.Path):
		// put admin method here
		return

    // NORMAL
	case r.Method == http.MethodPost && TagBaseRegex.MatchString(r.URL.Path):
		d.CreateTag(w, r)
		return
	case r.Method == http.MethodGet && TagBaseRegex.MatchString(r.URL.Path):
        d.GetTags(w, r)
		return
	case r.Method == http.MethodGet && TagByIdRegex.MatchString(r.URL.Path):
        d.GetTagById(w, r)
		return
	case r.Method == http.MethodGet && TagByNameRegex.MatchString(r.URL.Path):
        d.GetTagByName(w, r)
		return
	case r.Method == http.MethodPut && TagByIdRegex.MatchString(r.URL.Path):
		//d.UpdateUserKeyValuePair(w, r)
		return
	case r.Method == http.MethodDelete && TagByIdRegex.MatchString(r.URL.Path):
        d.DeleteTagById(w, r)
		return
	
	default:
		w.Write([]byte("Test response from tags"))
		return
	}
}

// FRONTEND ROUTES
func (d *TagsRoute) CreateTag(w http.ResponseWriter, r *http.Request) {
	var t models.Tag
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while decoding tag post from json", err)
		return
	}

	newTag, err := t.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while creating tag: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	jsonTag, err := json.Marshal(newTag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while marshalling tag to json: ", err)
		return
	}

	w.Write(jsonTag)
}

func (d *TagsRoute) GetTags(w http.ResponseWriter, r *http.Request) {
	tagSlice, err := models.ReadAllTags()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while getting tag list", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(tagSlice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while marshalling tag slice to json: ", err)
		return
	}

	w.Write(jsonData)
}

func (d *TagsRoute) GetTagById(w http.ResponseWriter, r *http.Request) {
	id_path := TagByIdRegex.FindStringSubmatch(r.URL.Path)
	tid, err := strconv.Atoi(id_path[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while getting id from path", err)
		return
	}

	t, err := models.ReadTagById(tid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while reading tag by id", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while marshalling tag to json while getting tag by id: ", err)
		return
	}

	w.Write(jsonData)
}

func (d *TagsRoute) GetTagByName(w http.ResponseWriter, r *http.Request) {
	id_path := TagByNameRegex.FindStringSubmatch(r.URL.Path)
	tagTitle := id_path[1]

	t, err := models.ReadTagByName(tagTitle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while reading tag by name", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while marshalling tag to json while getting tag by name: ", err)
		return
	}

	w.Write(jsonData)
}

func (d *TagsRoute) DeleteTagById(w http.ResponseWriter, r *http.Request) {
	id_path := TagByIdRegex.FindStringSubmatch(r.URL.Path)
	tid, err := strconv.Atoi(id_path[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while getting id from path during tag deletion", err)
		return
	}

	t, err := models.ReadTagById(tid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while reading tag by id during deletion", err)
		return
	}

	err = t.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while deleting tag", err)
		return
	}

	w.Write([]byte(fmt.Sprintf("Tag with id %d was deleted succesfully", tid)))
}
