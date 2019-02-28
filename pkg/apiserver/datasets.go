package apiserver

import (
	"net/http"

	"github.com/gorilla/mux"
	"minibox.ai/minibox/pkg/models"
	"minibox.ai/minibox/pkg/utils/tmpl"
)

type Datasets struct {
	*mux.Router
}

func NewDatasets() *Datasets {
	mux := &Datasets{
		Router: mux.NewRouter(),
	}
	mux.init()
	return mux
}

func (ds *Datasets) init() {
	ds.HandleFunc("/top", ds.top)
	ds.HandleFunc("/{dataset}", ds.detail)
}

func (ds *Datasets) top(w http.ResponseWriter, r *http.Request) {
	tmpl.RenderTemplate(w, "top.html.tmpl", nil)
}

func (ds *Datasets) detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dataset := vars["dataset"]
	d := models.GetDataset(dataset)
	if ds != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	tmpl.RenderTemplate(w, "dataset_top_detail.html.tmpl", d)
}
