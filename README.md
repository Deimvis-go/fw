# fw

fw (framework) — golang library that brings request and response abstractions.

## Install

```bash
go get github.com/Deimvis-go/fw
```

## Examples

### Define a typed request

```go
type V1UsersIdGETRequest struct {
    fw.RequestGET
    fw.RequestPathBound
    fw.RequestHeader[fw.JSONHeaderPreset]

    fw.RequestURI[struct {
        UserId string `uri:"id" validate:"required"`
    }]
    fw.RequestQuery[struct {
        IncludeDeleted bool `query:"include_deleted" form:"include_deleted"`
    }]
    fw.RequestNoBody
}

func (r *V1UsersIdGETRequest) Path() string {
    return fmt.Sprintf("/v1/users/%v", r.URI.UserId)
}
```

### Define a typed response

```go
type V1UsersIdGETResponse200 struct {
    fw.Response200
    fw.ResponseHeader[fw.JSONHeaderPreset]

    fw.ResponseBodyJSON[struct {
        Id   string `json:"id"`
        Name string `json:"name"`
    }]
}

// or use a pre-defined shortcut for responses without body:
type V1UsersIdGETResponse404 = fw.Response404WithJSONHeader
```

### Client-side flow

```go
ctx := context.Background()
client := http.DefaultClient

req := &V1UsersIdGETRequest{}
req.URI.UserId = "123"

// option 1: manual decoding
resp, close, err := fw.Do(ctx, client, req)
if err != nil {
    log.Fatal(err)
}
defer close()
if resp.Code() == 200 {
	var resp200 V1UsersIdGETResponse200
    err = fwresponse.move(resp &resp200) 
    // ...
} else {
    err := fmt.Errorf("unexpected response: %s", fwfmt.Response(resp))
    log.Fatal(err)
}

// option 2: use decoding option
gresp, close, err := us.a.Do(ctx, &req,
    fwhttp.WithMoveToFirstMatched(
      fwmatch.ByCode,
      &V1UsersIdGETResponse200,
      &V1UsersIdGETResponse404,
    ))
if err != nil {
    log.Fatal(err)
}
switch resp := gresp.(type) {
case *V1UsersIdGETResponse200:
    // ...
case *V1UsersIdGETResponse404:
    // ...
default:
    panic("bug: unreachable code")
```

### Server-side flow

#### Gin Framework

```go
// off-the-shelf solution for gin framework
// https://github.com/Deimvis-go/xgin
func GetUser(ctx context.Context, req &V1UsersIdGETRequest) fw.Resposne {
    user := db.GetUser(req.URI.UserId)
    if !user.HasValue() {
        return fwss.Resp404("no user(%s)", req.URI.UserId)
        // or &V1UsersIdGETResponse404{}
    }
    var resp V1UsersIdGETResponse200
    resp.Id = user.id
    resp.Name = user.Name
    return &resp
}

r.GET("/v1/users/:id", ginss.NewHandler(GetUser))
```
