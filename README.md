# crawler
Website crawler

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
