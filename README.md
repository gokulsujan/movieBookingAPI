# Movie Booking API
This project is written purely in Go Language. Gin (Http Web Frame Work) is used in this project. PostgreSQL Database is used to manage the data.
## Framework Used
Gin-Gonic: This whole project is built on Gin frame work. Its is a popular http web frame work.
```
go get -u github.com/gin-gonic/gin
```
## Database Used
PostgreSQL: PostgreSQL is a powerful, open source object-relational database. The data managment in this project is done using PostgreSQL. ORM tool named GORM is also been used to simplify the forms of queries for better understanding.

```
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```
## External Package Used
#### Razorpay
For Payment I have used the test case of Razorpay.
```
github.com/razorpay/razorpay-go
```