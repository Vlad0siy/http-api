package main

import (
  "net/http"
  "net/url"
  "time"
  "log"
  "encoding/json"
  "strconv"
  "fmt"
  "unicode"
  "strings"
)

type Response struct {
  Success bool
  Err string
  Value string
}

type MyError struct {
  When time.Time
  ErrMsg string
}

func (e MyError) Error() string {
  return fmt.Sprintf("%v %v", e.When, e.ErrMsg)
}

func RiseErr(Msg string) error {
  return MyError{
    time.Now(),
    Msg,
  }
}

func main() {
  handler := http.NewServeMux()
  handler.HandleFunc("/add", Logger(AddHandler))
  handler.HandleFunc("/sub", Logger(SubHandler))
  handler.HandleFunc("/mul", Logger(MulHandler))
  handler.HandleFunc("/div", Logger(DivHandler))

  server := http.Server {
    Addr: "localhost:6060",
    Handler: handler,
    ReadTimeout: 10 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  log.Fatal(server.ListenAndServe())
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
  int_a, int_b := 0, 0
  fl_a, fl_b := 0.0, 0.0
  var err error
  if IsFloat(r.URL.Query().Get("a")) || IsFloat(r.URL.Query().Get("b")) {
    fl_a, fl_b, err = GetFloatParams(r.URL.Query)
  } else {
    int_a, int_b, err = GetIntParams(r.URL.Query)
  }
  if err != nil {
    CreateResponse("0", w, err.Error())
    return
  }
  if int_a != 0 && int_b != 0 {
    CreateResponse(strconv.Itoa(int_a + int_b), w, "")
  } else {
    CreateResponse(strconv.FormatFloat(fl_a + fl_b, 'g', -1, 64), w, "")
  }
}

func SubHandler(w http.ResponseWriter, r *http.Request) {
  int_a, int_b := 0, 0
  fl_a, fl_b := 0.0, 0.0
  var err error
  if IsFloat(r.URL.Query().Get("a")) || IsFloat(r.URL.Query().Get("b")) {
    fl_a, fl_b, err = GetFloatParams(r.URL.Query)
  } else {
    int_a, int_b, err = GetIntParams(r.URL.Query)
  }
  if err != nil {
    CreateResponse("0", w, err.Error())
    return
  }
  if int_a != 0 && int_b != 0 {
    CreateResponse(strconv.Itoa(int_a - int_b), w, "")
  } else {
    CreateResponse(strconv.FormatFloat(fl_a - fl_b, 'g', -1, 64), w, "")
  }
}

func MulHandler(w http.ResponseWriter, r *http.Request) {
  int_a, int_b := 0, 0
  fl_a, fl_b := 0.0, 0.0
  var err error
  if IsFloat(r.URL.Query().Get("a")) || IsFloat(r.URL.Query().Get("b")) {
    fl_a, fl_b, err = GetFloatParams(r.URL.Query)
  } else {
    int_a, int_b, err = GetIntParams(r.URL.Query)
  }
  if err != nil {
    CreateResponse("0", w, err.Error())
    return
  }
  if int_a != 0 && int_b != 0 {
    CreateResponse(strconv.Itoa(int_a * int_b), w, "")
  } else {
    CreateResponse(strconv.FormatFloat(fl_a * fl_b, 'g', -1, 64), w, "")
  }
}

func DivHandler(w http.ResponseWriter, r *http.Request) {
  int_a, int_b := 0, 0
  fl_a, fl_b := 0.0, 0.0
  var err error
  if IsFloat(r.URL.Query().Get("a")) || IsFloat(r.URL.Query().Get("b")) {
    fl_a, fl_b, err = GetFloatParams(r.URL.Query)
  } else {
    int_a, int_b, err = GetIntParams(r.URL.Query)
  }
  if err != nil {
    CreateResponse("0", w, err.Error())
    return
  }
  if (int_b == 0 && int_a != 0) || (fl_b == 0.0 && fl_a != 0.0)  {
    log.Println("Integer divide by zero")
    err := RiseErr("Integer divide by zero")
    CreateResponse("0", w, err.Error())
  } else if int_a != 0 && int_b != 0 {
    CreateResponse(strconv.Itoa(int_a / int_b), w, "")
  } else {
    CreateResponse(strconv.FormatFloat(fl_a / fl_b, 'g', -1, 64), w, "")
  }
}

func GetFloatParams(q func() url.Values) (float64, float64, error) {
  if rune(q().Get("a")[0]) != '-' || rune(q().Get("b")[0]) != '-' {
    if !unicode.IsNumber(rune(q().Get("a")[0])) || !unicode.IsNumber(rune(q().Get("b")[0])) {
      return 0, 0, RiseErr("Not valid data")
    }
  }
  a, aerr := ToFloat(q().Get("a"))
  b, berr := ToFloat(q().Get("b"))
  fmt.Printf("%T %T\n", a, b)
  if fmt.Sprintf("%T", a) != "float64" || fmt.Sprintf("%T", b) != "float64" {
    return 0, 0, RiseErr("Not valid data")
  }
  if aerr != nil {
    return 0, 0, aerr
  } else if berr != nil {
    return 0, 0, berr
  }
  return a, b, nil
}

func GetIntParams(q func() url.Values) (int, int, error) {
  if rune(q().Get("a")[0]) != '-' || rune(q().Get("b")[0]) != '-' {
    if !unicode.IsNumber(rune(q().Get("a")[0])) || !unicode.IsNumber(rune(q().Get("b")[0])) {
      return 0, 0, RiseErr("Not valid data")
    }
  }
  a, aerr := strconv.Atoi(q().Get("a"))
  b, berr := strconv.Atoi(q().Get("b"))
  fmt.Printf("%T %T\n", a, b)
  if fmt.Sprintf("%T", a) != "int" || fmt.Sprintf("%T", b) != "int" {
    return 0, 0, RiseErr("Not valid data")
  }
  if aerr != nil {
    return 0, 0, aerr
  } else if berr != nil {
    return 0, 0, berr
  }
  return a, b, nil
}

func IsFloat(val string) bool {
  return strings.Index(val, ".") > 0 || strings.Index(val, ",") > 0
}

func ToFloat(val string) (float64, error) {
  if strings.Index(val, ",") > 0 {
    val = strings.Replace(val, ",", ".", 1)
  }
  a, err := strconv.ParseFloat(val, 64)
  return a, err
}

func CreateResponse(Result string, w http.ResponseWriter, Err string) {
  Resp := Response {
    Err == "",
    Err,
    "0",
  }

  Resp.Value = Result

  RespJson, err := json.Marshal(Resp)
  if err != nil {
    log.Fatal(err)
  }

  w.WriteHeader(http.StatusOK)
  w.Write(RespJson)
}

func Logger(Next http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    log.Printf("server [net/http] method [%s] connection from [%v]\n", r.Method, r.RemoteAddr)
    log.Printf("Params: a [%s]; b [%s]\n", r.URL.Query()["a"], r.URL.Query()["b"])
    Next.ServeHTTP(w, r)
  }
}
