package api

import (
    "fmt"
    "errors"
    "net/http"
    //"encoding/json"
    "github.com/gorilla/mux"

    "filestorage/view"
    "filestorage/model/bucket"
    "filestorage/model/object"
    "filestorage/instance/log"
)

//
type Server struct {
}

//
func New() Server {
    return Server{}
}

//
func (s *Server) Run(port string) error {
    router := mux.NewRouter()

    router.HandleFunc("/", wrapperReq(listBucket)).Methods("GET")
    router.HandleFunc("/", wrapperReq(createBucket)).Methods("POST")
    router.HandleFunc("/{name}", wrapperReq(getBucket)).Methods("GET")
    router.HandleFunc("/{name}", wrapperReq(deleteBucket)).Methods("DELETE")
    router.HandleFunc("/{name}/", wrapperReq(createObject)).Methods("POST")
    router.HandleFunc("/{name}/", wrapperReq(listObject)).Methods("GET")
    router.HandleFunc("/{name}/{file:.*}", wrapperReq(metaObject)).Methods("OPTIONS")
    router.HandleFunc("/{name}/{file:.*}", wrapperReq(fileObject)).Methods("GET")
    router.HandleFunc("/{name}/{file:.*}", wrapperReq(deleteObject)).Methods("DELETE")
    router.HandleFunc("/{name}/{file:.*}", wrapperReq(updateObject)).Methods("PUT")

    http.Handle("/", router)
    return http.ListenAndServe(port, nil)
}

//
func wrapperReq(handler http.HandlerFunc) http.HandlerFunc {
    return func(res http.ResponseWriter, req *http.Request) {
        log.Info("request", fmt.Sprintf("%+v",req))

        handler(res, req)

        log.Info("response", fmt.Sprintf("%+v",res))
        return
    }
}

// List bucket
func listBucket(res http.ResponseWriter, req *http.Request) {
    list, err := view.GetListBucket()
    if err != nil {
        toError(res, "LIST_BUCKET_NEW", err, http.StatusInternalServerError)
        return
    }

    toJSON(res, list, http.StatusOK)
    return
}

// Create bucket
func createBucket(res http.ResponseWriter, req *http.Request) {
    service := req.Header.Get("X-Basis-Service")
    log.Info("service", service)

    custom, err := fromBody(req)
    if err != nil {
        toError(res, "CREATE_NEW_CUSTOM", err, http.StatusInternalServerError)
        return
    }

    create, err := custom.NewCreateBucket(service)
    if err != nil {
        toError(res, "CREATE_BUCKET_NEW_CREATE", err, http.StatusBadRequest)
        return
    }

    err = create.ValidationCreateBucket()
    if err != nil {
        toError(res, "CREATE_BUCKET_VALIDATION", err, http.StatusBadRequest)
        return
    }

    if bucket.IsExistsBucketByName(create.Name) {
        toError(res, "CREATE_BUCKET_DUPLICATE", errors.New("Error create duplicate bucket"), http.StatusBadRequest)
        return
    }

    err = create.CreateInStorage()
    if err != nil {
        toError(res, "CREATE_BUCKET_CREATED", err, http.StatusInternalServerError)
        return
    }

    res.Header().Set("Location", create.GetLocation())
    res.WriteHeader(http.StatusCreated)
    return
}

// Get bucket by name
func getBucket(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    name := vars["name"]

    exists := bucket.IsExistsBucketByName(name)
    if !exists {
        toError(res, "GET_BUCKET_NOT_FOUND", errors.New("Get bucket not found"), http.StatusNotFound)
        return
    }

    bucket, err := bucket.GetBucketByName(name)
    if err != nil {
        toError(res, "GET_BUCKET_INTERNAL_ERROR", err, http.StatusInternalServerError)
        return
    }

    custom, err := view.NewCustomByBucket(bucket)
    if err != nil {
        toError(res, "GET_BUCKET_TO_CUSTOM", err, http.StatusInternalServerError)
        return
    }

    toJSON(res, custom, http.StatusOK)
    return
}

// Delete bucket
func deleteBucket(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    name := vars["name"]

    exists := bucket.IsExistsBucketByName(name)
    if !exists {
        toError(res, "DELETE_BUCKET_NOT_FOUND", errors.New("Delete bucket not found"), http.StatusNotFound)
        return
    }

    bucket, err := bucket.GetBucketByName(name)
    if err != nil {
        toError(res, "DELETE_BUCKET_INTERNAL_ERROR", err, http.StatusInternalServerError)
        return
    }

    deleted, err := bucket.DeleteBucket()
    if err != nil {
        toError(res, "DELETE_BUCKET_INTERNAL_ERROR", err, http.StatusInternalServerError)
        return
    } else if !deleted {
        toError(res, "DELETE_BUCKET_INTERNAL_ERROR", errors.New("Return delete false"), http.StatusInternalServerError)
        return
    }

    res.WriteHeader(http.StatusNoContent)
    return
}

// Create object
func createObject(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    name := vars["name"]

    exists := bucket.IsExistsBucketByName(name)
    if !exists {
        toError(res, "CREATE_OBJECT_BUCKET_NOT_FOUND", errors.New("Bucket not found for create object"), http.StatusNotFound)
        return
    }

    metadata := req.PostFormValue("metadata")
    log.Info("metadata", metadata)


/*
    obj, err := object.NewCreateObject()
    if err != nil {
        toError(res, "CREATE_OBJECT_NEW", err, http.StatusInternalServerError)
        return
    }
/*
    file, header, err := req.FormFile("datafile")
    if err != nil {
        toError(res, "CREATE_OBJECT_FORM_FILE", err, http.StatusInternalServerError)
        return
    }
    log.Info("filename", header.Filename, "Content-Type:", header.Header.Get("Content-Type"), "Size:", header.Size)

    obj.Bucket = name
    obj.ContentType = header.Header.Get("Content-Type")
    obj.Prefix = req.PostFormValue("prefix")
    obj.Name = req.PostFormValue("name")
    if obj.Name == "" {
        obj.Name = header.Filename
    }

    meta := req.PostFormValue("metadata")
    log.Info("metadata", meta)

    err = json.Unmarshal([]byte(meta), &obj.Metadata)
    if err != nil {
        toError(res, "CREATE_OBJECT_UNMARSHAL", err, http.StatusInternalServerError)
        return
    }

    err = obj.ValidationCreateObject()
    if err != nil {
        toError(res, "CREATE_OBJECT_VALIDATION", err, http.StatusBadRequest)
        return
    }

    exists = object.IsExistsObjectByKey(obj.Bucket, obj.GetCreateObjectKey())
    if  exists {
        toError(res, "CREATE_OBJECT_DUPLICATE", errors.New("Error create duplicate object"), http.StatusBadRequest)
        return
    }

    err = obj.CreateObjectInStorage(file, header.Size)
    if err != nil {
        toError(res, "CREATE_BUCKET_CREATED", err, http.StatusInternalServerError)
        return
    }

    res.Header().Set("Location", fmt.Sprintf("/%s/%s", obj.Bucket, obj.GetCreateObjectKey()))
*/
    res.WriteHeader(http.StatusCreated)
    return
}

// List object
func listObject(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    name := vars["name"]

    exists := bucket.IsExistsBucketByName(name)
    if !exists {
        toError(res, "LIST_OBJECT_BUCKET_NOT_FOUND", errors.New("Bucket not found for create object"), http.StatusNotFound)
        return
    }

    objects, err := object.GetListObjectByBucket(name, "")
    if err != nil {
        toError(res, "LIST_OBLECTS_BY_BUCKET_GET", err, http.StatusInternalServerError)
        return
    }

    toJSON(res, objects, http.StatusOK)
    return
}

// Get meta object
func metaObject(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    name := vars["name"]
    file := vars["file"]

    exists := bucket.IsExistsBucketByName(name)
    if !exists {
        toError(res, "META_OBJECT_BUCKET_NOT_FOUND", errors.New("Bucket not found for create object"), http.StatusNotFound)
        return
    }

    exists = object.IsExistsObjectByKey(name, file)
    if !exists {
        toError(res, "META_OBJECT_OBJECT_NOT_FOUND", errors.New("Object not found"), http.StatusNotFound)
        return
    }

    objects, err := object.GetListObjectVersion(name, file)
    if err != nil {
        toError(res, "META_OBLECT_NEW", err, http.StatusInternalServerError)
        return
    }

    toJSON(res, objects, http.StatusOK)
    return
}

// Delete Object
func deleteObject(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    name := vars["name"]
    file := vars["file"]

    exists := bucket.IsExistsBucketByName(name)
    if !exists {
        toError(res, "DELETE_OBJECT_BUCKET_NOT_FOUND", errors.New("Bucket not found"), http.StatusNotFound)
        return
    }

    exists = object.IsExistsObjectByKey(name, file)
    if !exists {
        toError(res, "DELETE_OBJECT_OBJECT_NOT_FOUND", errors.New("Object not found"), http.StatusNotFound)
        return
    }

    objects, err := object.GetListObjectVersion(name, file)
    if err != nil {
        toError(res, "DELETE_OBJECT_NEW", err, http.StatusInternalServerError)
        return
    }

    deleted, err := objects.DeleteAllObject()
    if err != nil {
        toError(res, "DELETE_OBJECT_INTERNAL_ERROR", err, http.StatusInternalServerError)
        return
    } else if !deleted {
        toError(res, "DELETE_OBJECT_INTERNAL_ERROR", errors.New("Return delete false"), http.StatusInternalServerError)
        return
    }

    res.WriteHeader(http.StatusNoContent)
    return
}

// Get file object
func fileObject(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    name := vars["name"]
    file := vars["file"]

    exists := bucket.IsExistsBucketByName(name)
    if !exists {
        toError(res, "FILE_OBJECT_BUCKET_NOT_FOUND", errors.New("Bucket not found for create object"), http.StatusNotFound)
        return
    }

    exists = object.IsExistsObjectByKey(name, file)
    if !exists {
        toError(res, "FILE_OBJECT_OBJECT_NOT_FOUND", errors.New("Object not found"), http.StatusNotFound)
        return
    }

    var obj *object.Object
    var err error

    versionId := req.URL.Query().Get("versionId")
    if versionId != "" {
        obj, err = object.NewObjectByVersion(name, file, versionId)
    } else {
        obj, err = object.NewObject(name, file)
    }
    if err != nil {
        toError(res, "FILE_OBLECT_NEW", err, http.StatusInternalServerError)
        return
    }

    data, err := obj.DataObject()
    if err != nil {
        toError(res, "FILE_OBLECT_DATA", err, http.StatusInternalServerError)
        return
    }

    res.Header().Set("Content-Type", obj.ContentType)
    res.WriteHeader(http.StatusOK)

    _, err = res.Write(data)
    if err != nil {
        toError(res, "FILE_OBLECT_DATA", err, http.StatusInternalServerError)
        return
    }

    return
}

// Update object
func updateObject(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    name := vars["name"]
    file := vars["file"]

    exists := bucket.IsExistsBucketByName(name)
    if !exists {
        toError(res, "UPDATE_OBJECT_BUCKET_NOT_FOUND", errors.New("Bucket not found for create object"), http.StatusNotFound)
        return
    }

    exists = object.IsExistsObjectByKey(name, file)
    if !exists {
        toError(res, "UPDATE_OBJECT_OBJECT_NOT_FOUND", errors.New("Object not found"), http.StatusNotFound)
        return
    }

    obj, err := object.NewObject(name, file)
    if err != nil {
        toError(res, "UPDATE_OBLECT_NEW", err, http.StatusInternalServerError)
        return
    }

    datafile, header, err := req.FormFile("datafile")
    if err != nil {
        toError(res, "CREATE_OBJECT_FORM_FILE", err, http.StatusInternalServerError)
        return
    }
    log.Info("filename", header.Filename, "Content-Type:", header.Header.Get("Content-Type"), "Size:", header.Size)
    log.Info("datafile", datafile)

    create, err := object.NewCreateObject()
    if err != nil {
        toError(res, "UPDATE_CREATE_OBLECT_NEW", err, http.StatusInternalServerError)
        return
    }

    create.Bucket = name
    create.Name = obj.Name
    create.Prefix = obj.Prefix
    create.Metadata = obj.Metadata
    create.ContentType = header.Header.Get("Content-Type")

    err = create.ValidationCreateObject()
    if err != nil {
        toError(res, "UPDATE_OBJECT_VALIDATION", err, http.StatusBadRequest)
        return
    }

    err = create.CreateObjectInStorage(datafile, header.Size)
    if err != nil {
        toError(res, "UPDATE_BUCKET_CREATED", err, http.StatusInternalServerError)
        return
    }

    res.WriteHeader(http.StatusOK)
    return
}

