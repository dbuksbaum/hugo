---
title: transform.Unmarshal
description: Parses serialized data and returns a map or an array. Supports CSV, JSON, TOML, YAML, and XML.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [unmarshal]
    returnType: any
    signatures: ['transform.Unmarshal [OPTIONS] INPUT']
aliases: [/functions/transform.unmarshal]
---

The input can be a string or a [resource](g).

## Unmarshal a string

```go-html-template
{{ $string := `
title: Les Misérables
author: Victor Hugo
`}}

{{ $book := unmarshal $string }}
{{ $book.title }} → Les Misérables
{{ $book.author }} → Victor Hugo
```

## Unmarshal a resource

Use the `transform.Unmarshal` function with global, page, and remote resources.

### Global resource

A global resource is a file within the `assets` directory, or within any directory mounted to the `assets` directory.

```text
assets/
└── data/
    └── books.json
```

```go-html-template
{{ $data := dict }}
{{ $path := "data/books.json" }}
{{ with resources.Get $path }}
  {{ with . | transform.Unmarshal }}
    {{ $data = . }}
  {{ end }}
{{ else }}
  {{ errorf "Unable to get global resource %q" $path }}
{{ end }}

{{ range where $data "author" "Victor Hugo" }}
  {{ .title }} → Les Misérables
{{ end }}
```

### Page resource

A page resource is a file within a [page bundle].

```text
content/
├── post/
│   └── book-reviews/
│       ├── books.json
│       └── index.md
└── _index.md
```

```go-html-template
{{ $data := dict }}
{{ $path := "books.json" }}
{{ with .Resources.Get $path }}
  {{ with . | transform.Unmarshal }}
    {{ $data = . }}
  {{ end }}
{{ else }}
  {{ errorf "Unable to get page resource %q" $path }}
{{ end }}

{{ range where $data "author" "Victor Hugo" }}
  {{ .title }} → Les Misérables
{{ end }}
```

### Remote resource

A remote resource is a file on a remote server, accessible via HTTP or HTTPS.

```go-html-template
{{ $data := dict }}
{{ $url := "https://example.org/books.json" }}
{{ with try (resources.GetRemote $url) }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else with .Value }}
    {{ $data = . | transform.Unmarshal }}
  {{ else }}
    {{ errorf "Unable to get remote resource %q" $url }}
  {{ end }}
{{ end }}

{{ range where $data "author" "Victor Hugo" }}
  {{ .title }} → Les Misérables
{{ end }}
```

> [!note]
> When retrieving remote data, a misconfigured server may send a response header with an incorrect [Content-Type]. For example, the server may set the Content-Type header to `application/octet-stream` instead of `application/json`.
>
> In these cases, pass the resource `Content` through the `transform.Unmarshal` function instead of passing the resource itself. For example, in the above, do this instead:
>
> `{{ $data = .Content | transform.Unmarshal }}`

## Working with CSV

### Options

When unmarshaling a CSV file, provide an optional map of options.

delimiter
: (`string`) The delimiter used. Default is `,`.

comment
: (`string`) The comment character used in the CSV. If set, lines beginning with the comment character without preceding whitespace are ignored.

lazyQuotes
: {{< new-in 0.122.0 />}}
: (`bool`) Whether to allow a quote in an unquoted field, or to allow a non-doubled quote in a quoted field. Default is `false`.

targetType
: {{< new-in 0.146.7 />}}
: (`string`) The target data type, either `slice` or `map`. Default is `slice`.

### Examples

The examples below use this CSV file:

```csv
"name","type","breed","age"
"Spot","dog","Collie",3
"Rover","dog","Boxer",5
"Felix","cat","Calico",7
```

To render an HTML table from a CSV file:

```go-html-template
{{ $data := slice }}
{{ $file := "pets.csv" }}
{{ with or (.Resources.Get $file) (resources.Get $file) }}
  {{ $opts := dict "targetType" "slice" }}
  {{ $data = transform.Unmarshal $opts . }}
{{ end }}

{{ with $data }}
  <table>
    <thead>
      <tr>
        {{ range index . 0 }}
          <th>{{ . }}</th>
        {{ end }}
      </tr>
    </thead>
    <tbody>
      {{ range . | after 1 }}
        <tr>
          {{ range . }}
            <td>{{ . }}</td>
          {{ end }}
        </tr>
      {{ end }}
    </tbody>
  </table>
{{ end }}
```

To extract a subset of the data, or to sort the data, unmarshal to a map instead of a slice:

```go-html-template
{{ $data := slice }}
{{ $file := "pets.csv" }}
{{ with or (.Resources.Get $file) (resources.Get $file) }}
  {{ $opts := dict "targetType" "map" }}
  {{ $data = transform.Unmarshal $opts . }}
{{ end }}

{{ with sort (where $data "type" "dog") "name" "asc" }}
  <table>
    <thead>
      <tr>
        <th>name</th>
        <th>type</th>
        <th>breed</th>
        <th>age</th>
      </tr>
    </thead>
    <tbody>
      {{ range . }}
        <tr>
          <td>{{ .name }}</td>
          <td>{{ .type }}</td>
          <td>{{ .breed }}</td>
          <td>{{ .age }}</td>
        </tr>
      {{ end }}
    </tbody>
  </table>
{{ end }}
```

## Working with XML

When unmarshaling an XML file, do not include the root node when accessing data. For example, after unmarshaling the RSS feed below, access the feed title with `$data.channel.title`.

```xml
<?xml version="1.0" encoding="utf-8" standalone="yes"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>Books on Example Site</title>
    <link>https://example.org/books/</link>
    <description>Recent content in Books on Example Site</description>
    <language>en-US</language>
    <atom:link href="https://example.org/books/index.xml" rel="self" type="application/rss+xml" />
    <item>
      <title>The Hunchback of Notre Dame</title>
      <description>Written by Victor Hugo</description>
      <link>https://example.org/books/the-hunchback-of-notre-dame/</link>
      <pubDate>Mon, 09 Oct 2023 09:27:12 -0700</pubDate>
      <guid>https://example.org/books/the-hunchback-of-notre-dame/</guid>
    </item>
    <item>
      <title>Les Misérables</title>
      <description>Written by Victor Hugo</description>
      <link>https://example.org/books/les-miserables/</link>
      <pubDate>Mon, 09 Oct 2023 09:27:11 -0700</pubDate>
      <guid>https://example.org/books/les-miserables/</guid>
    </item>
  </channel>
</rss>
```

Get the remote data:

```go-html-template
{{ $data := dict }}
{{ $url := "https://example.org/books/index.xml" }}
{{ with try (resources.GetRemote $url) }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else with .Value }}
    {{ $data = . | transform.Unmarshal }}
  {{ else }}
    {{ errorf "Unable to get remote resource %q" $url }}
  {{ end }}
{{ end }}
```

Inspect the data structure:

```go-html-template
<pre>{{ debug.Dump $data }}</pre>
```

List the book titles:

```go-html-template
{{ with $data.channel.item }}
  <ul>
    {{ range . }}
      <li>{{ .title }}</li>
    {{ end }}
  </ul>
{{ end }}
```

Hugo renders this to:

```html
<ul>
  <li>The Hunchback of Notre Dame</li>
  <li>Les Misérables</li>
</ul>
```

### XML attributes and namespaces

Let's add a `lang` attribute to the `title` nodes of our RSS feed, and a namespaced node for the ISBN number:

```xml
<?xml version="1.0" encoding="utf-8" standalone="yes"?>
<rss version="2.0"
  xmlns:atom="http://www.w3.org/2005/Atom"
  xmlns:isbn="http://schemas.isbn.org/ns/1999/basic.dtd"
>
  <channel>
    <title>Books on Example Site</title>
    <link>https://example.org/books/</link>
    <description>Recent content in Books on Example Site</description>
    <language>en-US</language>
    <atom:link href="https://example.org/books/index.xml" rel="self" type="application/rss+xml" />
    <item>
      <title lang="en">The Hunchback of Notre Dame</title>
      <description>Written by Victor Hugo</description>
      <isbn:number>9780140443530</isbn:number>
      <link>https://example.org/books/the-hunchback-of-notre-dame/</link>
      <pubDate>Mon, 09 Oct 2023 09:27:12 -0700</pubDate>
      <guid>https://example.org/books/the-hunchback-of-notre-dame/</guid>
    </item>
    <item>
      <title lang="fr">Les Misérables</title>
      <description>Written by Victor Hugo</description>
      <isbn:number>9780451419439</isbn:number>
      <link>https://example.org/books/les-miserables/</link>
      <pubDate>Mon, 09 Oct 2023 09:27:11 -0700</pubDate>
      <guid>https://example.org/books/les-miserables/</guid>
    </item>
  </channel>
</rss>
```

After retrieving the remote data, inspect the data structure:

```go-html-template
<pre>{{ debug.Dump $data }}</pre>
```

Each item node looks like this:

```json
{
  "description": "Written by Victor Hugo",
  "guid": "https://example.org/books/the-hunchback-of-notre-dame/",
  "link": "https://example.org/books/the-hunchback-of-notre-dame/",
  "number": "9780140443530",
  "pubDate": "Mon, 09 Oct 2023 09:27:12 -0700",
  "title": {
    "#text": "The Hunchback of Notre Dame",
    "-lang": "en"
  }
}
```

The title keys do not begin with an underscore or a letter---they are not valid [identifiers](g). Use the [`index`] function to access the values:

```go-html-template
{{ with $data.channel.item }}
  <ul>
    {{ range . }}
      {{ $title := index .title "#text" }}
      {{ $lang := index .title "-lang" }}
      {{ $ISBN := .number }}
      <li>{{ $title }} ({{ $lang }}) {{ $ISBN }}</li>
    {{ end }}
  </ul>
{{ end }}
```

Hugo renders this to:

```html
<ul>
  <li>The Hunchback of Notre Dame (en) 9780140443530</li>
  <li>Les Misérables (fr) 9780451419439</li>
</ul>
```

[`index`]: /functions/collections/indexfunction/
[Content-Type]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type
[page bundle]: /content-management/page-bundles/
