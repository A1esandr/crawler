# crawler
Website crawler

## Usage

### Prerequisites
* Go 1.16

### Example 
```
links, err := crawler.New().Run(url)
if err != nil {
  fmt.Println(err)
}

for _, link := range links {
  fmt.Println(link)
}
```
