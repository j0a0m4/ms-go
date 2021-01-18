package infra

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

// Products slice stores Product in-memory
type Products []*Product

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	_, idx, err := findProduct(id)

	if err != nil {
		return err
	}

	p.ID = id
	p.UpdatedOn = time.Now().UTC().String()
	productList[idx] = p

	return nil
}

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, errors.New("Product not found")
}

func getNextID() int {
	return len(productList) + 1
}

// ToJSON encodes Products to JSON
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "L47S",
		CratedOn:    time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   "",
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee withour milk",
		Price:       1.99,
		SKU:         "32PS",
		CratedOn:    time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
		DeletedOn:   "",
	},
}
