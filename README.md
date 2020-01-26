# gRPC Örnek REST API

Oluşturulan repository go lang ögrenimi ve eklentisi olan **gRPC** ile **protobuff** formatı oluşturmak ile beraber gRPC ile microservice HTTP/REST endpoint, middleware.

Mimari yapıyı hazırlarken [golang-standarts/project-layout](https://github.com/golang-standards/project-layout) structure'ını kullanmaktayız, standartın dışına çıkmamak adına.Böylelikle proje yönetimimiz daha kolay erişebilirlik anlamında bize rahatlık sağlatacaktır.

# İçerik

Bu repository'nin mimarisinde bulunanlar aşağıda 3 madde de özetlenmiştir.

-   ***gRPC CRUD service*** ve ***client*** oluşturmak
   	> RAW SQL Query'ler yazılacaktır ilerleyen zamanlarda istenilirse  ORM & Model 3. parti go frameworkleri ile bütünleşebilir.
   	
-   ***HTTP-REST*** endpointini ***gRPC*** servisine bağlamak
-   ***HTTP-REST*** ve ***gRPC*** servisine ***Middleware*** (logging/tracing) eklemek

## Referanslar

- [*Language Guide (proto3)*](https://developers.google.com/protocol-buffers/docs/proto3)
