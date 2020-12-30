# Hệ thống CRUD người dùng đơn giản (Gin/Gonic, JWT, Gorm, Postgres)
#Giới thiệu

Hệ thống hiện thực các API bên dưới:

1. POST /users  (1)  Body: {"first_name":"your_surname", "last_name":"your name", "email": "your_email", "password":"XXXXXXXX" }   
2. POST /login (2)   Body: {"email": "your_email", "password":"XXXXXXXX" } - Return {"access_token": xxxxxxxx, "refresh_token": xxxxxxxxx} ()
3. POST /logout (3)  
4. GET  /profile?access_token (4)  
5. PUT  /update?access_token (5) Body: format xxx-form-encoded  [Key: email     Value: your_new_email ]
6. DELETE /delete/{:user_id} (6)
7. GET /users/{:user_id}?access_token (7)
8. GET /users?access_token (8)

On Postman go to:

Authentication tab
Select type: Bearer Token
Paste in your Token you received after send request login

#Các bước chạy chương trình

Trên Windows 10:

1. Tải Postgres (phiên bản 9.6) trên trang chủ : https://www.postgresql.org/download/ , cài đặt Postgres trên Windows, setup user="postgres", password="123456", giữ port mặc định sẵn (5432). 
2. Tải pgAdmin4 trên trang chủ : https://www.pgadmin.org/download/, sau đó cài đặt trên Windows
3. Tải, giải nén file zip, cài đặt Redis trên Windows, open redis-server để khởi động Redis, db mặc định chạy trên địa chỉ IP: 127.0.0.1:6379 
4. (Optional) Tải và cài đặt Postman trên desktop để test các API.
5. Truy cập pgAdmin4, create database mang tên "customer".
6. Clone folder về máy, đặt folder vào $GOPATH/src, sau đó chạy lệnh "go run main.go" trên Terminal.

#Testing API.

Header {
   1. OPTIONS localhost:8888 (1) | (2) | (3) | (4) | (5) | (6) | (7) | (8)
   2. Access-Control-Allow-Origin: *
   3. Access-Control-Allow-Credentials: true
   4. Access-Control-Request-Method: GET | POST | PUT | DELETE
   5. Access-Control-Request-Headers: Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With
}


