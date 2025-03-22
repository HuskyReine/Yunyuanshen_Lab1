@echo off

echo Stopping gomall-cart...
docker stop gomall-cart

echo Stopping gomall-checkout...
docker stop gomall-checkout

echo Stopping gomall-email...
docker stop gomall-email

echo Stopping gomall-frontend...
docker stop gomall-frontend

echo Stopping gomall-order...
docker stop gomall-order

echo Stopping gomall-payment...
docker stop gomall-payment

echo Stopping gomall-product...
docker stop gomall-product

echo Stopping gomall-user...
docker stop gomall-user

echo All containers stopped.
pause