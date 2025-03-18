// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"context"
	"strconv"

	// "errors"
	// "fmt"
	// "strconv"

	// "github.com/cloudwego/biz-demo/gomall/app/checkout/infra/mq"
	// "github.com/cloudwego/biz-demo/gomall/app/checkout/infra/rpc"
	"github.com/cloudwego/biz-demo/gomall/app/checkout/infra/mq"
	"github.com/cloudwego/biz-demo/gomall/app/checkout/infra/rpc"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart"
	checkout "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/checkout"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/email"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/payment"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type CheckoutService struct {
	ctx context.Context
} // NewCheckoutService new CheckoutService
func NewCheckoutService(ctx context.Context) *CheckoutService {
	return &CheckoutService{ctx: ctx}
}

/*
	Run

1. get cart
2. calculate cart
3. create order
4. empty cart
5. pay
6. change order result
7. finish
*/
func (s *CheckoutService) Run(req *checkout.CheckoutReq) (resp *checkout.CheckoutResp, err error) {

	// TODO 1.get cart (使用RPC调用Cart服务以获得购物车信息)
	cartResult, err := rpc.CartClient.GetCart(s.ctx, &cart.GetCartReq{UserId: req.UserId})
	if err != nil {
		return nil, kerrors.NewGRPCBizStatusError(114514, "failed to get cart: "+err.Error())
	}
	if cartResult == nil || cartResult.Cart.Items == nil {
		return nil, kerrors.NewGRPCBizStatusError(114514, "your cart is empty!")
	}

	// TODO 2.calc cart（根据第1步的购物车信息，计算总价和订单项信息）
	var (
		total float32
		oi    []*order.OrderItem
	)

	for _, cartItem := range cartResult.Cart.Items {

		// 获取商品信息，RPC调用Product服务
		productResp, err := rpc.ProductClient.GetProduct(s.ctx, &product.GetProductReq{
			Id: cartItem.ProductId,
		})
		if err != nil {
			return nil, kerrors.NewGRPCBizStatusError(114514, "failed to get product info: "+err.Error())
		}
		if productResp.Product == nil {
			continue
		}

		// 计算单项成本
		cost := productResp.Product.Price * float32(cartItem.Quantity)
		total += cost

		// 添加订单项
		oi = append(oi, &order.OrderItem{
			Item: &cart.CartItem{
				ProductId: cartItem.ProductId,
				Quantity:  cartItem.Quantity,
			},
			Cost: cost,
		})

	}

	// TODO 3.create order（根据第1步和第2步的信息，创建order.PlaceOrderReq，并使用RPC调用Order服务创建订单）
	zipcode, _ := strconv.ParseInt(req.Address.ZipCode, 10, 32)
	orderResp, err := rpc.OrderClient.PlaceOrder(s.ctx, &order.PlaceOrderReq{
		UserId: req.UserId,
		Email:  req.Email,
		Address: &order.Address{
			StreetAddress: req.Address.StreetAddress,
			City:          req.Address.City,
			State:         req.Address.State,
			Country:       req.Address.Country,
			ZipCode:       int32(zipcode),
		},
		OrderItems: oi,
	})

	if err != nil {
		return nil, kerrors.NewGRPCBizStatusError(114514, "failed to place order: "+err.Error())
	}
	if orderResp == nil || orderResp.Order == nil {
		return nil, kerrors.NewGRPCBizStatusError(114514, "invalid order response")
	}
	orderId := orderResp.Order.OrderId

	// TODO 4.empty cart（使用RPC调用Cart服务清空购物车）
	_, err = rpc.CartClient.EmptyCart(s.ctx, &cart.EmptyCartReq{UserId: req.UserId})
	if err != nil {
		klog.Warnf("failed to empty cart for user %d: %v", req.UserId, err)
	}

	// TODO 5.pay（使用RPC调用Payment服务进行支付）
	paymentResult, err := rpc.PaymentClient.Charge(s.ctx, &payment.ChargeReq{
		UserId:  req.UserId,
		OrderId: orderId,
		Amount:  total,
		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          req.CreditCard.CreditCardNumber,
			CreditCardCvv:             req.CreditCard.CreditCardCvv,
			CreditCardExpirationMonth: req.CreditCard.CreditCardExpirationMonth,
			CreditCardExpirationYear:  req.CreditCard.CreditCardExpirationYear,
		},
	})
	if err != nil {
		return nil, kerrors.NewGRPCBizStatusError(114514, "payment failed: "+err.Error())
	}
	if paymentResult == nil || paymentResult.TransactionId == "" {
		return nil, kerrors.NewGRPCBizStatusError(114514, "invalid payment response")
	}

	// TODO 6.send email（使用MQ发送邮件通知）
	data, _ := proto.Marshal(&email.EmailReq{
		From:        "from@example.com",
		To:          req.Email,
		ContentType: "text/plain",
		Subject:     "You just created an order in CloudWeGo shop",
		Content:     "You just created an order in CloudWeGo shop",
	})
	msg := &nats.Msg{Subject: "email", Data: data}
	err = mq.Nc.PublishMsg(msg)
	if err != nil {
		klog.Error(err.Error())
	}
	klog.Info(paymentResult)

	// TODO 7.finish（返回订单ID和支付结果）
	resp = &checkout.CheckoutResp{
		OrderId:       orderId,
		TransactionId: paymentResult.TransactionId,
	}
	return

}
