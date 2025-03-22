REM 需要手动按8次ctrl+c

@echo off

echo Creating gomall-cart...
docker run -it --network=gomall --name=gomall-cart gomall-cart /app/cart/output/bin/cart
docker stop gomall-cart

echo Creating gomall-checkout...
docker run -it --network=gomall --name=gomall-checkout gomall-checkout /app/checkout/output/bin/checkout
docker stop gomall-checkout

echo Creating gomall-email
docker run -it --network=gomall --name=gomall-email gomall-email /app/email/output/bin/email
docker stop gomall-email

echo Creating gomall-frontend
docker run -it --network=gomall -p 8180:8180 --name=gomall-frontend gomall-frontend /app/frontend/output/bin/frontend
docker stop gomall-frontend

echo Creating gomall-order
docker run -it --network=gomall --name=gomall-order gomall-order /app/order/output/bin/order
docker stop gomall-order

echo Creating gomall-payment
docker run -it --network=gomall --name=gomall-payment gomall-payment /app/payment/output/bin/payment
docker stop gomall-payment

echo Creating gomall-product
docker run -it --network=gomall --name=gomall-product gomall-product /app/product/output/bin/product
docker stop gomall-product

echo Creating gomall-user
docker run -it --network=gomall --name=gomall-user gomall-user /app/user/output/bin/user
docker stop gomall-user

echo All containers created.