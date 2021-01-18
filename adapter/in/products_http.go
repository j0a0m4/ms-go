package in

import (
	"errors"
	"log"
	"ms-go/infra"
	"net/http"
	"regexp"
	"strconv"
)

// ProductsHTTP is a RESTful Controller
type ProductsHTTP struct {
	l *log.Logger
}

// NewProductsHTTP instanciates a new HttpProducts controller
func NewProductsHTTP(l *log.Logger) *ProductsHTTP {
	return &ProductsHTTP{l}
}

func (p *ProductsHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.l.Printf("%s %s ", r.Method, r.RequestURI)

	switch r.Method {
	case http.MethodGet:
		p.getProducts(w, r)
	case http.MethodPost:
		p.addProduct(w, r)
	case http.MethodPut:
		id, err := p.parseURL(r.URL.Path)
		if err != nil {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
		}
		p.updateProduct(id, w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (p *ProductsHTTP) parseURL(path string) (int, error) {
	p.l.Printf("parseUrl; start; path=%s\n", path)

	regex := regexp.MustCompile(`/([0-9]+)`)
	group := regex.FindAllStringSubmatch(path, -1)

	if len(group) != 1 || len(group[0]) != 2 {
		p.l.Panicln("parseUrl; end; error; message=Invalid URI")
		return -1, errors.New("Invalid URI")
	}

	id, err := strconv.Atoi(group[0][1])
	p.l.Printf("parseUrl; end; success; id=%v\n", id)
	return id, err
}

func (p *ProductsHTTP) getProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("getProducts; start;")

	products := infra.GetProducts()
	err := products.ToJSON(w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		p.l.Panicln("getProducts; end; error; exception=Unable to marshall JSON")
		return
	}

	p.l.Printf("getProducts; end; success; listSize=%d\n", len(products))
}

func (p *ProductsHTTP) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("addProducts; start;")

	product := &infra.Product{}
	err := product.FromJSON(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		p.l.Panicln("addProducts; end; error; exception=Unable to marshall JSON")
		return
	}

	infra.AddProduct(product)
	w.WriteHeader(http.StatusCreated)
	p.l.Printf("addProducts; end; success; body=%#v\n", product)
}

func (p *ProductsHTTP) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	p.l.Printf("updateProduct; start; id: %v\n", id)

	product := &infra.Product{}
	err := product.FromJSON(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		p.l.Panicln("updateProduct; end; error; exception=Unable to marshall JSON")
		return
	}

	err = infra.UpdateProduct(id, product)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		p.l.Panicln("updateProduct; end; error; exception=Product not found")
		return
	}

	p.l.Printf("updateProduct; end; success;\n")
}
