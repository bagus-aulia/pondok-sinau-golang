# Perpustakaan Pondok Sinau - Go

Aplikasi ini menyediakan API untuk Front End Perpustakaan Pondok Sinau. Aplikasi ini menggunakan library seperti

### 1. [Gin dari gin-gonic](https://github.com/gin-gonic/gin)

Digunakan untuk mempermudah route dan proses input data. Untuk menginstall library, menggunakan command :
<pre>
$ go get -u github.com/gin-gonic/gin
</pre>

### 2. [JWT Go dari dgrijalva](https://github.com/dgrijalva/jwt-go)

Digunakan untuk membuat token hak akses guna membatasi public user mengakses data. Untuk menginstall library, menggunakan command :
<pre>
$ go get github.com/dgrijalva/jwt-go
</pre>

[Documentation](https://godoc.org/github.com/dgrijalva/jwt-go)

### 3. [Gotenv by subosito](https://github.com/subosito/gotenv)

Used to load environment variables from `.env` or `io.Reader` in Go. To add this library to your PC, run this command in terminal:
<pre>
$ go get github.com/subosito/gotenv
</pre>

### 4. [Gorm by jinzhu](https://github.com/jinzhu/gorm)

ORM library for golang. Used to connect with database. Supported database are MySQL, PostgreSQL, Sqlite3 and SQL Server. To add this library to your PC, run this command in terminal:
<pre>
$ go get github.com/jinzhu/gorm
</pre>

[Documentation](https://gorm.io/docs/index.html)

### 5. [Gocialite by danilopolani](https://github.com/jinzhu/gorm)

Gocialite used to manage social oAuth authentication. This library available for Amazon, Asana, Bitbucket, Facebook, Foursquare, Github, Google, LinkedIn or Slack account. To add this library to your PC, run this command in terminal:
<pre>
$ go get gopkg.in/danilopolani/gocialite.v1
</pre>

### 6. [Slug by gosimple](https://github.com/gosimple/slug)

Package slug generate slug from unicode string, URL-friendly slugify with multiple languages support. Used to generate username if its blank. To add this library to your PC, run this command in terminal:
<pre>
$ go get github.com/gosimple/slug
</pre>

[Documentation](https://godoc.org/github.com/gosimple/slug)

2020 - Bagus Aulia Al Ilhami
