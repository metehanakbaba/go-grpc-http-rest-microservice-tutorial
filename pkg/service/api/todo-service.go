package api

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/metehanakbaba/go-grpc-http-rest-microservice-tutorial/pkg/api/api"
)

// toDoServiceServer'e Proto'da ki TodoServiceServer interfacesini implement edelim
type toDoServiceServer struct {
	db *sql.DB
}

// Servisimizi olusturalim
func NewToDoServiceServer(db *sql.DB) api.ToDoServiceServer {
	return &toDoServiceServer{db: db}
}

// SQL Database baglanti fonksiyonu
func (s *toDoServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "veritabani baglanti hatasi "+err.Error())
	}
	return c, nil
}

// Olustur
func (s *toDoServiceServer) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	// SQL Baglantimizi pool'dan cekelim
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	reminder, err := ptypes.Timestamp(req.ToDo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder fieldi yanlis formatda"+err.Error())
	}

	// CreateRequest den gelen datalari sqle into edelim
	res, err := c.ExecContext(ctx, "INSERT INTO ToDo(`Title`, `Description`, `Reminder`) VALUES(?, ?, ?)",
		req.ToDo.Title, req.ToDo.Description, reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Create(*context, *api.CreateRequest) SQL Hata:  "+err.Error())
	}

	// Olusturulan son id yi cekelim
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "olusturulan todo verileri retrieve sorunu"+err.Error())
	}

	return &api.CreateResponse{
		Id: id,
	}, nil
}

// Okuma
func (s *toDoServiceServer) Read(ctx context.Context, req *api.ReadRequest) (*api.ReadResponse, error) {
	// SQL Baglantimizi pool'dan cekelim
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT `ID`, `Title`, `Description`, `Reminder` FROM ToDo WHERE `ID`=?",
		req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Read(*context, *api.ReadRequest) SQL Hata"+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "Veri alinamadi"+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo ID='%d' bulunamadi",
			req.Id))
	}

	// SQL Verilerini (api.ToDo) Interfaceni
	var td api.ToDo
	var reminder time.Time
	if err := rows.Scan(&td.Id, &td.Title, &td.Description, &reminder); err != nil {
		return nil, status.Error(codes.Unknown, "SQL verileri interface'e aktarilamadi"+err.Error())
	}

	td.Reminder, err = ptypes.TimestampProto(reminder)

	if err != nil {
		return nil, status.Error(codes.Unknown, "reminder fieldi yanlis formatda "+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("Birden fazla ID='%d' bulundu",
			req.Id))
	}

	return &api.ReadResponse{
		ToDo: &td,
	}, nil

}

// Guncelleme
func (s *toDoServiceServer) Update(ctx context.Context, req *api.UpdateRequest) (*api.UpdateResponse, error) {
	// SQL Baglantimizi pool'dan cekelim
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, "UPDATE ToDo SET `Title`=?, `Description`=?, `Reminder`=? WHERE `ID`=?",
		req.ToDo.Title,
		req.ToDo.Description,
		time.Unix(req.ToDo.Reminder.Seconds, int64(req.ToDo.Reminder.Nanos)),
		req.ToDo.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Update(*context, *api.UpdateRequest) SQL Hata"+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "reminder fieldi yanlis formatda "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo ID='%d' bulunamadi",
			req.ToDo.Id))
	}

	return &api.UpdateResponse{
		Updated: rows,
	}, nil
}

// Silme
func (s *toDoServiceServer) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {

	// SQL Baglantimizi pool'dan cekelim
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// ToDo silme
	res, err := c.ExecContext(ctx, "DELETE FROM ToDo WHERE `ID`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Delete(*context, *api.DeleteRequest) SQL Hata"+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "Silinen ToDo verisi retrieve edilemedi"+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo ID='%d' bulunamadi",
			req.Id))
	}

	return &api.DeleteResponse{
		Deleted: rows,
	}, nil
}

// Tumunu Okuma
func (s *toDoServiceServer) ReadAll(ctx context.Context, req *api.ReadAllRequest) (*api.ReadAllResponse, error) {

	// SQL Baglantimizi pool'dan cekelim
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// ToDo listesi getir
	rows, err := c.QueryContext(ctx, "SELECT `ID`, `Title`, `Description`, `Reminder` FROM ToDo")
	if err != nil {
		return nil, status.Error(codes.Unknown, "ToDo Select Hatasi "+err.Error())
	}
	defer rows.Close()

	var reminder time.Time
	list := []*api.ToDo{}

	for rows.Next() {
		td := new(api.ToDo)
		if err := rows.Scan(&td.Id, &td.Title, &td.Description, &reminder); err != nil {
			return nil, status.Error(codes.Unknown, "ToDo satirindan veri alinamadi"+err.Error())
		}
		td.Reminder, err = ptypes.TimestampProto(reminder)
		if err != nil {
			return nil, status.Error(codes.Unknown, "Reminder sutunu formati yanlis veri turunde bind edilmis"+err.Error())
		}
		list = append(list, td)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "ToDo tumunu oku methodundan veri alinamadi "+err.Error())
	}

	return &api.ReadAllResponse{
		ToDos: list,
	}, nil
}
