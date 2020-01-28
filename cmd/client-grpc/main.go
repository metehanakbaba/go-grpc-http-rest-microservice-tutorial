package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"

	"github.com/metehanakbaba/go-grpc-http-rest-microservice-tutorial/pkg/api/api"
)

func main() {
	// flagler
	address := flag.String("server", "", "gRPC Server host:port")
	flag.Parse()

	// gRPC Serverini kuralim
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Baglanti saglanamadi: %v", err)
	}
	defer conn.Close()

	//
	c := api.NewToDoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Reminder sütununu yeniden override ediyoruz (gRPC Protobuf pTypes.Timestamp) olarak
	t := time.Now().In(time.UTC)            // Tarih olusturma
	reminder, _ := ptypes.TimestampProto(t) // (t time.Time) (*tspb.Timestamp, error)
	pfx := t.Format(time.RFC3339Nano)       // tarih formatini yazdirabilecegimiz sekilde RFC3339Nano = 2006-01-02T15:04:05.999999999Z07:00

	req1 := api.CreateRequest{
		ToDo: &api.ToDo{
			Title:       "Baslik (" + pfx + ")",
			Description: "Icerik (" + pfx + ")",
			Reminder:    reminder,
		},
	}

	res1, err := c.Create(ctx, &req1)
	if err != nil {
		log.Fatalf("Olusturma hatasi: %v", err)
	}
	log.Printf("Olusturulan sonuc: <%+v>\n\n", res1)

	id := res1.Id

	req2 := api.ReadRequest{
		Id: id,
	}
	res2, err := c.Read(ctx, &req2)

	if err != nil {
		log.Fatalf("Okuma hatasi: %v", err)
	}
	log.Printf("Okuma sonucu: <%+v>\n\n", res2)

	req3 := api.UpdateRequest{
		ToDo: &api.ToDo{
			Id:          res2.ToDo.Id,
			Title:       res2.ToDo.Title,
			Description: res2.ToDo.Description + " + guncellendi",
			Reminder:    res2.ToDo.Reminder,
		},
	}
	res3, err := c.Update(ctx, &req3)
	if err != nil {
		log.Fatalf("Guncelleme hastası: %v", err)
	}
	log.Printf("Guncelleme sonucu: <%+v>\n\n", res3)

	req4 := api.ReadAllRequest{}
	res4, err := c.ReadAll(ctx, &req4)
	if err != nil {
		log.Fatalf("Tumunu okuma hatasi: %v", err)
	}
	log.Printf("Okuma sonucu: <%+v>\n\n", res4)

	req5 := api.DeleteRequest{
		Id: id,
	}
	res5, err := c.Delete(ctx, &req5)
	if err != nil {
		log.Fatalf("Silme hatasi: %v", err)
	}
	log.Printf("Silme sonucu: <%+v>\n\n", res5)
}
