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
	"fmt"

	// "github.com/cloudwego/biz-demo/gomall/app/order/biz/dal/mysql"
	// "github.com/cloudwego/biz-demo/gomall/app/order/biz/model"
	"github.com/cloudwego/biz-demo/gomall/app/order/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/order/biz/model"
	order "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	// "github.com/google/uuid"
	// "gorm.io/gorm"
)

type PlaceOrderService struct {
	ctx context.Context
} // NewPlaceOrderService new PlaceOrderService
func NewPlaceOrderService(ctx context.Context) *PlaceOrderService {
	return &PlaceOrderService{ctx: ctx}
}

// Run create note info
func (s *PlaceOrderService) Run(req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {

	// TODO 请实现PlaceOrder的业务逻辑，插入数据到数据库中的order表和order_item表，生成一个随机的uuid作为订单号
	// 可以参考其他服务的源代码实现这个函数

	// klog.Warnf("PlaceOrderService.Run not implemented")
	// resp = &order.PlaceOrderResp{
	// 	Order: &order.OrderResult{
	// 		OrderId: "1145141919810",
	// 	},
	// }
	// return

	// 校验订单项是否为空
	if len(req.OrderItems) == 0 {
		err = kerrors.NewBizStatusError(191981, "items is empty")
		return
	}

	// 生成唯一的订单号UUID
	orderUUID, err := uuid.NewUUID()
	if err != nil {
		err = fmt.Errorf("failed to generate UUID: %w", err)
		return
	}
	orderId := orderUUID.String()

	// 开启事务，插入订单和订单项数据
	// 使用gorm的事务，确保订单及其订单项插入数据库时保持一致性
	err = mysql.DB.Transaction(func(tx *gorm.DB) error {

		// 构建订单数据
		o := &model.Order{
			OrderId: orderId,
			UserId:  req.UserId,
			Consignee: model.Consignee{
				Email: req.Email,
			},
		}

		// 填充地址信息
		if req.Address != nil {
			a := req.Address
			o.Consignee.StreetAddress = a.StreetAddress
			o.Consignee.City = a.City
			o.Consignee.State = a.State
			o.Consignee.Country = a.Country
		}

		// 插入订单order
		if err := tx.Create(o).Error; err != nil {
			return err
		}

		// 构建订单项数据
		items := make([]model.OrderItem, 0, len(req.OrderItems))
		for _, v := range req.OrderItems {
			items = append(items, model.OrderItem{
				OrderIdRefer: orderId,
				ProductId:    v.Item.ProductId,
				Quantity:     v.Item.Quantity,
				Cost:         v.Cost,
			})
		}

		// 插入订单项order_item
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
		return nil

	})

	// 事务失败时返回业务错误
	if err != nil {
		return nil, kerrors.NewBizStatusError(191981, "place order failed: "+err.Error())
	}

	// 成功返回订单ID
	resp = &order.PlaceOrderResp{
		Order: &order.OrderResult{
			OrderId: orderId,
		},
	}
	return

}
