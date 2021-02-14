package mock_ecommerce

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb "productinfo/client/ecommerce"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	wrapper "github.com/golang/protobuf/ptypes/wrappers"
)

// rpcMsg implements the gomock.Matcher interface
type rpcMsg struct {
	msg proto.Message
}

func (r *rpcMsg) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m, r.msg)
}

func (r *rpcMsg) String() string {
	return fmt.Sprintf("is %s", r.msg)
}

//---------------------------------------------------------
// 코드 7-2 부분
//---------------------------------------------------------
func TestAddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mocklProdInfoClient := NewMockProductInfoClient(ctrl)
	//...
	name := "Sumsung S21"
	description := "Samsung Galaxy S21 is the latest smart phone, launched in Jan. 2021"
	price := float32(700.0)

	req := &pb.Product{Name: name, Description: description, Price: price}

	mocklProdInfoClient.
		EXPECT().AddProduct(gomock.Any(), &rpcMsg{msg: req}).
		Return(&wrapper.StringValue{Value: "ABC123" + name}, nil)

	testAddProduct(t, mocklProdInfoClient)
}

func testAddProduct(t *testing.T, client pb.ProductInfoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//...
	name := "Sumsung S21"
	description := "Samsung Galaxy S21 is the latest smart phone, launched in Jan. 2021"
	price := float32(700.0)

	r, err := client.AddProduct(ctx, &pb.Product{Name: name,
		Description: description, Price: price})
	// 테스트와 응답 검증

	if err != nil || r.GetValue() != "ABC123Sumsung S21" {
		t.Errorf("mocking failed")
	}
	t.Log("Reply : ", r.GetValue())
}

//---------------------------------------------------------
