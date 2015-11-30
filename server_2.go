package main

import (
"github.com/julienschmidt/httprouter"
    "strconv"
    "net/http"
    
    "bytes"
"fmt" 
    "log"
)

var hmap map[int]string


func listing(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
    if len(hmap) == 0 {
        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(200)
        fmt.Fprintf(writer, "{}")
    } else {
        var buff bytes.Buffer
        buff.WriteString("{\n[\n")
        for k, v := range hmap {
            str := `{
                        "key " : "` + strconv.Itoa(k) + `",
                        "value" : "` + v + `"
                    },` + "\n"
            buff.WriteString(str)
        }

        jstr := buff.String()
        jstr = jstr[:len(jstr)-2]
        jstr = jstr + "\n]\n}"

        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(400)
        fmt.Fprintf(writer, "%s\n", jstr)
    }
}

func get(writer http.ResponseWriter, _ *http.Request, parameter httprouter.Params) {
    k, err := strconv.Atoi(parameter.ByName("key_id"))

    if err != nil {
        writer.Header().Set("Content-Type", "plain/text")
        writer.WriteHeader(400)
        fmt.Fprintf(writer, "Key needs to be integer\n")
    } else {
        v, x := hmap[k]
        if x {
            writer.Header().Set("Content-Type", "application/json")
            writer.WriteHeader(400)

            jstr := `{
                "key" : "` + strconv.Itoa(k) + `",
                "value" : "` + v + `"
            }`

            fmt.Fprintf(writer, "%s\n", jstr)
        } else {
            writer.Header().Set("Content-Type", "plain/text")
            writer.WriteHeader(400)
            fmt.Fprintf(writer, "invalid request\n")
        }
    }
}

func addition(writer http.ResponseWriter, _ *http.Request, parameter httprouter.Params) {
    v := parameter.ByName("value")
    k, err := strconv.Atoi(parameter.ByName("key_id"))

    if err != nil {
        writer.Header().Set("Content-Type", "plain/text")
        writer.WriteHeader(400)
        fmt.Fprintf(writer, "Key should be integer\n")
    } else {
        hmap[k] = v
        writer.Header().Set("Content-Type", "plain/text")
        writer.WriteHeader(200)
        fmt.Fprintf(writer, "Added {key:%x, value:%s}", k, v)
    }
}

func main() {
    hmap = make(map[int]string)

    r := httprouter.New()
    r.PUT("/keys/:k_id/:v", addition)
    r.GET("/keys/:k_id", get)
    r.GET("/keys", listing)

    log.Println("Server listening 4001")
    http.ListenAndServe(":4001", r)
}
