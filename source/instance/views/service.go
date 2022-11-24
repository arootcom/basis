package views

import (
)

//
type Service struct {
    Service string                      `json:"service"`
    Objects []*Object                   `json:"objects"`
}

//
type Object struct {
    Base string                         `json:"base"`
    View View                           `json:"view"`
}

//
type View struct {
    Type string                         `json:"type"`
    Attributes []Attribute              `json:"attributes"`
}

//
type Attribute struct {
    Name string                         `json:"name"`
    Filled Value                        `json:"filled"`
    Tags bool                           `json:"tags"`
    Validation []interface{}            `json:"validation"`
}

//
type Value interface{
}

