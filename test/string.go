package tester

/*
Constants Generator
Help to auto-generated correct constants for lists
Martin Prestone (c) 2024
github.com/monopolly

Struct options
type=string    Type of const list. Ex type=string, int by default
#swift         Create swift enum
#index         Create int index anyway
#up            Uppercase for golang vars
name=Status    Named var. Ex name=Category

Fields options
name{}         Add new name for fields. Ex name{newjsonname}
go{}           Add new golang var name for fields. Ex go{ArtList}
title{}        Add title for fields. Ex title{News}
desc{}         Add desc for fields. Ex desc{New movie is good}
#              Add lists for fields. Ex: id int //#readonly #must...
list{}         Create a.List() for field. Ex: list{colors, text}
iota=100       Add iota counter
*/
type Category string

const (
	CHECKBOX = Category("checkbox") //title{Checkbox}
	NUMBER   = Category("number")   //title{Number}
	TEXT     = Category("text")     //title{Text}
	MARKDOWN = Category("markdown") //title{Markdown}
	PRICE    = Category("price")    //title{Price} list{currency, price}
	COUNTRY  = Category("country")  //title{Country}
	CITY     = Category("city")     //title{City} list{geoid, country, state, city}
	IMAGE    = Category("image")    //title{Image}
	FILE     = Category("file")     //title{File}
	LINK     = Category("link")     //title{Link}
)

func (a Category) String() string {
	switch a {
	case CHECKBOX:
		return "checkbox"
	case NUMBER:
		return "number"
	case TEXT:
		return "text"
	case MARKDOWN:
		return "markdown"
	case PRICE:
		return "price"
	case COUNTRY:
		return "country"
	case CITY:
		return "city"
	case IMAGE:
		return "image"
	case FILE:
		return "file"
	case LINK:
		return "link"
	default:
		return ""
	}
}

func (a Category) Title() string {
	switch a {
	case CHECKBOX:
		return "Checkbox"
	case NUMBER:
		return "Number"
	case TEXT:
		return "Text"
	case MARKDOWN:
		return "Markdown"
	case PRICE:
		return "Price"
	case COUNTRY:
		return "Country"
	case CITY:
		return "City"
	case IMAGE:
		return "Image"
	case FILE:
		return "File"
	case LINK:
		return "Link"
	default:
		return ""
	}
}

func (a Category) Desc() string {
	switch a {
	default:
		return ""
	}
}

func ValidCategoryKey(v string) bool {
	switch Category(v) {
	case CHECKBOX, NUMBER, TEXT, MARKDOWN, PRICE, COUNTRY, CITY, IMAGE, FILE, LINK:
		return true
	default:
		return false
	}
}

func CategoryIndexes() []Category {
	return []Category{CHECKBOX, NUMBER, TEXT, MARKDOWN, PRICE, COUNTRY, CITY, IMAGE, FILE, LINK}
}

func (a Category) List() []string {
	switch a {
	case PRICE:
		return []string{"currency", "price"}
	case CITY:
		return []string{"geoid", "country", "state", "city"}
	default:
		return nil
	}
}
