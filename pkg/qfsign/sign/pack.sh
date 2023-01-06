#/bin/sh

if [ -d release ];then
    rm -rf release
fi

echo "start pack license tool."

output=release/license_tool
mkdir -p $output

\cp -rf pem $output
\cp -rf help.md $output


CGO_ENABLED=0 go build -o $output/sign -mod=vendor

echo "pack license tool ok."
