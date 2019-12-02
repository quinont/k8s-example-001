package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fwd := &forwarder{"words.api.svc.cluster.local", 8080}
	http.Handle("/words/", http.StripPrefix("/words", fwd))
	fwd2 := &forwarder{"app2.app2.svc.cluster.local", 5000}
	http.Handle("/app2/", http.StripPrefix("/app2", fwd2))
	http.Handle("/", http.FileServer(http.Dir("static")))

	fmt.Println("Listening on port 80")
	http.ListenAndServe(":80", nil)
}

type forwarder struct {
	host string
	port int
}

func (f *forwarder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addrs, err := net.LookupHost(f.host)
	if err != nil {
		log.Println("Error", err)
		http.Error(w, err.Error(), 500)
		return
	}

	log.Printf("%s %d available ips: %v", r.URL.Path, len(addrs), addrs)
	ip := addrs[rand.Intn(len(addrs))]
	log.Printf("%s I choose %s", r.URL.Path, ip)

	url := fmt.Sprintf("http://%s:%d%s", ip, f.port, r.URL.Path)
	log.Printf("%s Calling %s", r.URL.Path, url)

	test := r.Header.Get("test")

	if err = copy(url, ip, w, test); err != nil {
		log.Println("Error", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func copy(url, ip string, w http.ResponseWriter, test string) error {
	var resp *http.Response
	var err_connexion error

	if test == "si" {
		timeout := time.Duration(5 * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		request.Header.Add("test", "si")
		resp, err_connexion = client.Do(request)
	} else {
		resp, err_connexion = http.Get(url)
	}

	if err_connexion != nil {
		return err_connexion
	}

	w.Header().Set("source", ip)

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	return err
}
