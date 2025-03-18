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

	// "github.com/cloudwego/biz-demo/gomall/app/order/biz/dal/mysql"
	// "github.com/cloudwego/biz-demo/gomall/app/order/biz/model"
	"github.com/cloudwego/biz-demo/gomall/app/order/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/order/biz/model"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart"
	order "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type ListOrderService struct {
	ctx context.Context
} // NewListOrderService new ListOrderService
func NewListOrderService(ctx context.Context) *ListOrderService {
	return &ListOrderService{ctx: ctx}
}

// Run create note info
func (s *ListOrderService) Run(req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {

	// TODO 请实现ListOrder的业务逻辑，从数据库中的order表和order_item表中查询数据
	// 可以参考其他服务的源代码实现这个函数

	// klog.Warnf("ListOrderService.Run not implemented")
	// defaultOrder := &order.Order{
	// 	OrderId:      "1145141919810",
	// 	UserId:       114514,
	// 	UserCurrency: "JPY",
	// 	Address: &order.Address{
	// 		StreetAddress: "114514",
	// 		City:          "114514",
	// 		State:         "114514",
	// 		Country:       "JP",
	// 		ZipCode:       114514,
	// 	},
	// 	OrderItems: []*order.OrderItem{},
	// 	CreatedAt:  1145141919,
	// }
	// resp = &order.ListOrderResp{
	// 	Orders: []*order.Order{defaultOrder},
	// }
	// return

	// 从数据库中查询用户的订单及订单项数据
	list, err := model.ListOrder(mysql.DB, s.ctx, req.UserId)
	if err != nil {
		return nil, kerrors.NewBizStatusError(514114, err.Error())
	}

	// 订单
	var orders []*order.Order
	for _, v := range list {

		// 预分配订单项切片容量
		items := make([]*order.OrderItem, 0, len(v.OrderItems))
		for _, oi := range v.OrderItems {
			items = append(items, &order.OrderItem{
				Item: &cart.CartItem{
					ProductId: oi.ProductId,
					Quantity:  oi.Quantity,
				},
				Cost: oi.Cost,
			})
		}

		// 构建订单对象，填充收件人及地址信息
		orders = append(orders, &order.Order{
			CreatedAt: int32(v.CreatedAt.Unix()),
			OrderId:   v.OrderId,
			UserId:    v.UserId,
			Email:     v.Consignee.Email,
			Address: &order.Address{
				StreetAddress: v.Consignee.StreetAddress,
				Country:       v.Consignee.Country,
				City:          v.Consignee.City,
				State:         v.Consignee.State,
				ZipCode:       v.Consignee.ZipCode,
			},
			OrderItems: items,
		})
	}

	// 返回订单列表响应
	resp = &order.ListOrderResp{
		Orders: orders,
	}
	return

}
