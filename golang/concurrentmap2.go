package main

import ( "fmt" ; "sync" )

type DNS struct {
    m map[string]string
    lock sync.Mutex
}

func MakeDNS() *DNS {
    dns := new(DNS)
    dns.m = make(map[string]string)
    return dns
}

func (dns *DNS) GetAndSetOnce(suf string) {
    dns.lock.Lock() 
    defer dns.lock.Unlock()
    dns.m["X"]=dns.m["X"]+suf
}

func (dns *DNS) GetAndSet(suf string) {
    for i:=0; i<10; i++ { dns.GetAndSetOnce(suf) }
    c<-0
}

var c = make(chan int)

func main() {     
    dns := MakeDNS()
    go dns.GetAndSet("1") ; go dns.GetAndSet("2") ; go dns.GetAndSet("3")
    <-c ; <-c ; <-c // wait for the three goroutines to end
    fmt.Println(dns.m["X"])     
}