package main_test

import (
  "testing"
  "time"
  "net/http"
  "main"
)

func TestMain(t *testing.T) {
  handler := http.NewServeMux()
  handler.HandleFunc("/add", Logger(AddHandler))
  handler.HandleFunc("/sub", Logger(SubHandler))
  handler.HandleFunc("/mul", Logger(MulHandler))
  handler.HandleFunc("/div", Logger(DivHandler))

  TestTable := []struct {
    params []string
  } {
    {
      params: []string{"/add?a=-1&b=1", "/add?a=0&b=0", "/add?a=8&b=12", "/add?a=2.5&b=3.1", "/add?a=2,5&b=3,1", "/add?a=a&b=b",
                       "/sub?a=-1&b=1", "/sub?a=0&b=0", "/sub?a=8&b=12", "/sub?a=2.5&b=3.1", "/sub?a=2,5&b=3,1", "/sub?a=a&b=b",
                       "/mul?a=-1&b=1", "/mul?a=0&b=0", "/mul?a=8&b=12", "/mul?a=2.5&b=3.1", "/mul?a=2,5&b=3,1", "/mul?a=a&b=b",
                       "/div?a=-1&b=1", "/div?a=0&b=0", "/div?a=8&b=12", "/div?a=2.5&b=3.1", "/div?a=2,5&b=3,1", "/div?a=a&b=b"},
    },
  }
  server := http.Server {
    Addr: "localhost:8080",
    Handler: handler,
    ReadTimeout: 10 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  for _, TestParams := range TestTable {
    for _, Param := range TestParams.params {
      server.Addr += Param
      server.ListenAndServe()
    }
  }
}
