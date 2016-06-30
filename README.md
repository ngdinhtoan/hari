# Hari [![Build Status](https://travis-ci.org/ngdinhtoan/hari.svg)](https://travis-ci.org/ngdinhtoan/hari)

*Generate GO struct from JSON string*

### Why Hari?

When you make a client for RESTful service in GO, usually you have to define some structures to parse JSON response.
It's boring and takes time. Hari will help you on that job.

But, please be aware that generated struct may not meet your requirement. You have to review it before using.

### How to use

Install Hari by running the following command

    go get github.com/ngdinhtoan/hari

Put JSON string into a file within `.json` extension, the file name will be used to name struct.

Example file `product.json` has content

```json
{
    "id": 1,
    "name": "A green door",
    "price": 12.50,
    "active": true,
    "tags": ["home", "green"],
    "category": {
        "id": 2,
        "name": "Home"
    }
}
```

Run Hari command

    hari --input-dir=[path/to/dir]

then it will generate `product.go` file in the same directory, and its content is

```go
package main

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	Active   bool     `json:"active"`
	Category Category `json:"category"`
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Price    float64  `json:"price"`
	Tags     []string `json:"tags"`
}
```

### License

Hari is licensed under the [MIT License](https://github.com/ngdinhtoan/hari/blob/master/LICENSE).
