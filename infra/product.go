package infra

import (
	"encoding/json"
	"io"
	"ms-go/domain"
	"time"
)

// Product defines the structure for an API product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CratedOn    string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// ToJSON encodes Product to JSON
func (p *Product) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

// FromJSON decodes JSON to Product
func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	err := d.Decode(p)
	
	p.CratedOn = time.Now().UTC().String()
	p.UpdatedOn = time.Now().UTC().String()

	return err
}

// ToEntity maps the data object to a domain object
func (p *Product) ToEntity() *domain.Product {
	return &domain.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		SKU:         p.SKU,
	}
}
