# The best templating engine
Instead of explaining why it's the best, here's an example:

```html
<!DOCTYPE html>
<html>
    <head>
        <title>Example</title>
    </head>
    <body>
        <| posts := []string{ "Foo", "Bar", "Baz" } |>
        <ul>
            <| for _, post := range posts { |>
                <li><|:w post |></li>
            <| } |>
        </ul>
    </body>
</html>
```
after processing, this becomes exactly what you would expect:
```go
w(`<!DOCTYPE html>
<html>
    <head>
        <title>Example</title>
    </head>
    <body>
`);
posts := []string{ "Foo", "Bar", "Baz" };
w(`        <ul>
`);
for _, post := range posts {
w(`                <li>"); w(post); w("</li>
`);
}
w(`        </ul>
    </body>
</html>
`)
```
Golang!

Yes, this is an html to golang to html processor. Genius, I know. That's why I didn't come up with it. Zozin did!


