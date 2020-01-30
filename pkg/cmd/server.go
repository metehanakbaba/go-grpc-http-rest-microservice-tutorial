package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"

	// mysql driver gerekirse orm tercihi yapilir suanlik
	_ "github.com/go-sql-driver/mysql"

	"github.com/metehanakbaba/go-grpc-http-rest-microservice-tutorial/pkg/protocol/grpc"
	"github.com/metehanakbaba/go-grpc-http-rest-microservice-tutorial/pkg/protocol/rest"
	"github.com/metehanakbaba/go-grpc-http-rest-microservice-tutorial/pkg/service/api"
)

// Sunucu ayarları (istenilirse dot ENV tercihi yapılabilir default value'ler icin)
type Config struct {
	// gRPC server sunucu baslatma parametreleri
	// gRPC TCP port
	GRPCPort string

	// HTTP PORT
	HTTPPort string

	// MySQL sunucu baslatma parametreleri

	MySQLHost     string
	MySQLUser     string
	MySQLPassword string
	MySQLSchema   string
	MySQLParams   string
}

func RunServer() error {
	ctx := context.Background()

	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "Bind edilecek gRPC port")
	flag.StringVar(&cfg.HTTPPort, "http-port", "8080", "Bind edilecek HTTP/REST port")
	flag.StringVar(&cfg.MySQLHost, "mysql-host", "", "MySQL sunucu")
	flag.StringVar(&cfg.MySQLUser, "mysql-user", "", "MySQL kullanici")
	flag.StringVar(&cfg.MySQLPassword, "mysql-password", "", "MySQL sifre")
	flag.StringVar(&cfg.MySQLSchema, "mysql-schema", "", "MySQL veritabani adi")
	flag.StringVar(&cfg.MySQLParams, "mysql-params", "", "MySQL baglanti parametreleri ()")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 { // Cfg port kontrol etme
		return fmt.Errorf("Hatali gRPC server portu: '%s'", cfg.GRPCPort)
	}

	// Tarih formati icin otomatik parametre bind edelim
	if len(cfg.MySQLParams) > 0 { // Sql parametresi kontrol etme
		cfg.MySQLParams += "&"
	}
	cfg.MySQLParams += "parseTime=true"

	// sprintF ile datalari parse edip baglanti komutunu tanimlayalim
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		cfg.MySQLUser,
		cfg.MySQLPassword,
		cfg.MySQLHost,
		cfg.MySQLSchema,
		cfg.MySQLParams)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("Database Baglanti Hatasi : %v", err)
	}
	defer db.Close()

	v1API := api.NewToDoServiceServer(db)

	// HTTP gateway calistiralim
	go func() {
		err := rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
		log.Fatalf("HTTP/REST gateway baslatamadi: %v", err)
	}()

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
