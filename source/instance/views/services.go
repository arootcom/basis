package views

import (
    "fmt"
    "errors"
    "woodchuck/instance/log"
)

//
type Services struct {
    Items map[string]*Service
}

//
func (s *Services) GetViewByType(service string, base string, name string) (*View, error) {
    log.Debug("start", "GetViewByType: service =", service, ", base =", base, ", name =", name)

    srv, exists := s.Items[service]
    if !exists {
        err := errors.New("Service not fount")
        log.Error("error", err)
        return nil, err
    }

    var ret *View
    for _, obj := range srv.Objects {
        if obj.Base == base && obj.View.Type == name {
            ret = &obj.View
            break
        }
    }
    if ret == nil {
        err := errors.New("Object not found")
        log.Error("error", err)
        return nil, err
    }

    log.Debug("success", "GetViewByType:", fmt.Sprintf("%+v", ret))
    return ret, nil
}

