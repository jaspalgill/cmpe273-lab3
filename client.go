package main

import (
    "fmt"
    "log"
    "net/http"
    "strconv"
    "math/big"
    "crypto/md5"
    "encoding/hex"
    "io/ioutil"
    "github.com/julienschmidt/httprouter"
)

func tst() {
    
    hf := md5.New()
    vs := []string {"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
    for a, v := range vs {
                hf.Write([]byte(strconv.Itoa(a)))
		hv := big.NewInt(0)
		pnum := hv.Int64()
		h_str := hex.EncodeToString(hf.Sum(nil))
		hv.SetString(h_str, 16)

		
		if pnum < 0 {
			pnum *= -1
		}
        pnum = pnum % 3 + 3000

        url := "http://localhost:" + strconv.FormatInt(pnum, 10) + "/keys/" + strconv.Itoa(a+1) + "/" + v
        request, err := http.NewRequest("PUT", url, nil)
        if err != nil {
            log.Fatal(err)
        }
        _, err = http.DefaultClient.Do(request)
        if err != nil {
            log.Fatal(err)
        }
	}
}


func addition(writer http.ResponseWriter, _ *http.Request, parameter httprouter.Params) {
   v := parameter.ByName("v")
    k := parameter.ByName("key_id")
    
    pnum := hash(k)

    url := "http://localhost:" + strconv.FormatInt(pnum, 10) + "/keys/" + k + "/" + v
    request, err := http.NewRequest("PUT", url, nil)

    if err != nil {
        log.Fatal(err)
    }

    response, err := http.DefaultClient.Do(request)
    if err != nil {
        log.Fatal(err)
    }

    b, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
	defer response.Body.Close()

    if err != nil {
        writer.Header().Set("Content-Type", "plain/text")
        writer.WriteHeader(400)
        fmt.Fprintf(writer, "Key should be integer\n")
    } else {
        writer.Header().Set("Content-Type", "plain/text")
        writer.WriteHeader(400)
        fmt.Fprintf(writer, string(b) + ":localhost:%d", pnum)
    }
}
func hashing(k_id string) int64 {
    hf := md5.New()
    hv := big.NewInt(0)
    hf.Write([]byte(k_id))
    h_str := hex.EncodeToString(hf.Sum(nil))
    hv.SetString(h_str, 16)
    pnum := hv.Int64()
    if pnum < 0 {
        pnum *= -1
    }
    return pnum % 3 + 3000
}
func get(writer http.ResponseWriter, _ *http.Request, parameter httprouter.Params) {
    k := parameter.ByName("key_id")
    pnum := hash(k)

    url := "http://localhost:" + strconv.FormatInt(pnum, 10) + "/keys/" + k
    response, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    b, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
	defer response.Body.Close()

    if err != nil {
        writer.Header().Set("Content-Type", "plain/text")
        writer.WriteHeader(400)
        fmt.Fprintf(writer, "Key hould be integer\n")
    } else {
        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(200)
        fmt.Fprintf(writer, string(b))
    }
}


func main() {
    r := httprouter.New()
    r.PUT("/keys/:k_id/:v", addition)
    r.GET("/keys/:k_id", get)

    log.Println("Server listening on 8080")
    http.ListenAndServe(":8080", r)
}
