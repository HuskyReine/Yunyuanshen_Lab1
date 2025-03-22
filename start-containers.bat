@echo off

echo Starting gomall-cart...
docker start gomall-cart

echo Starting gomall-checkout...
docker start gomall-checkout

echo Starting gomall-email...
docker start gomall-email

echo Starting gomall-frontend...
docker start gomall-frontend

echo Starting gomall-order...
docker start gomall-order

echo Starting gomall-payment...
docker start gomall-payment

echo Starting gomall-product...
docker start gomall-product

echo Starting gomall-user...
docker start gomall-user

echo All containers started.
pause