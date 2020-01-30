# gRPC Örnek REST API

Oluşturulan repository go lang ögrenimi ve eklentisi olan **gRPC** ile **protobuff** formatı oluşturmak ile beraber gRPC ile microservice HTTP/REST endpoint, middleware.

Mimari yapıyı hazırlarken [golang-standarts/project-layout](https://github.com/golang-standards/project-layout) structure'ını kullanmaktayız, standartın dışına çıkmamak adına.Böylelikle proje yönetimimiz daha kolay erişebilirlik anlamında bize rahatlık sağlatacaktır.

Golang yazarlari [code review](https://github.com/golang/go/wiki/CodeReviewComments) yapilirken onerdikleri prosüdürleri uygulamaya gayret gosterildi.

Proje icerisinde [Degisken isimleri](https://github.com/golang/go/wiki/CodeReviewComments#variable-names) req (Request),res (Response),err (Error) şeklinde kısa tanımlamalar verildi.Bütün oluşabilecek hatalar handle edildi oluşturduğumuz servisin herhangi bir bug tracking yaşamaması adına.


## İçerik

Bu repository'nin mimarisinde bulunanlar aşağıda 3 madde de özetlenmiştir.

-   ***gRPC CRUD service*** ve ***client*** oluşturmak
   	> RAW SQL Query'ler yazılacaktır ilerleyen zamanlarda istenilirse  ORM & Model 3. parti go frameworkleri ile bütünleşebilir.
   	
-   ***HTTP-REST*** endpointini ***gRPC*** servisine bağlamak
-   ***HTTP-REST*** ve ***gRPC*** servisine ***Middleware*** (logging/tracing) eklemek

### Build & Kurulum

Veritabani icin MySQL sunucusuna ToDo tablosu ekleyelim

```
CREATE TABLE `ToDo` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `Title` varchar(200) DEFAULT NULL,
  `Description` varchar(1024) DEFAULT NULL,
  `Reminder` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `ID_UNIQUE` (`ID`)
);
```

- #### Build komutlari
    ```
    sh third_party/protoc-gen.sh
    cd cmd/server && go build .
    cd ../client-grpc && go build .
    ```
- #### gRPC Server
    ```
     ./cmd/server:
      -grpc-port string
            Bind edilecek gRPC port
      -mysql-host string
            MySQL sunucu
      -mysql-params string
            MySQL baglanti parametreleri ()
      -mysql-password string
            MySQL sifre
      -mysql-schema string
            MySQL veritabani adi
      -mysql-user string
            MySQL kullanici

    ```
    ##### Örnek komut satırı
    
    ```
     ./<binary> -grpc-port=<GRPC PORTU> -http-port=<HTTP_PORTU> -mysql-host=<MYSQL_SERVER_IP>:<PORT> -mysql-user=<KULLANICI ADI> -mysql-password=<SIFRE> --mysql-schema=<VERITABANI>
    ```
- #### gRPC Client
    ```
     ./cmd/client-grpc:
             -server string
                   gRPC Server host:port
    ```
    ##### Örnek komut satırı
    
    ```
     ./<binary>  -server=<GRPC_SERVER_IP>:<GRPC_PORT>
    ```  

### Önemli Go Modul Versiyonlari

- GO SDK 1.13.5
- GRPC 2.2.8


## Referanslar

- [*Language Guide (proto3)*](https://developers.google.com/protocol-buffers/docs/proto3)
