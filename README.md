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

Yes, this is an plaintext to golang to plaintext processor. Genius, I know. That's why I didn't come up with it. Zozin did!

You can actually use this with any plaintext document, or in theory actually any file but, I'm not sure what kind of edge cases that might introduce.

For a more complicated example, see the [section 7 manpage](./man/te.7.tet)

