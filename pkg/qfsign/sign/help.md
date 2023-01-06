1、生成签名许可
./sign  -mode rsa -private-key-cert  ./pem/rsa_private_key.pem
生成的签名值保存在/opt/openresty/nginx/.license/casb.license，默认使用本机某一个mac地址生成license
如果需要对某个mac地址生成，可以使用 -data 传入mac地址值

2、验签
./sign  -mode rsa -verify  -public-key-cert ./pem/rsa_public_key.pem
