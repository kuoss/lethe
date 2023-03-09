cd $(dirname $0)
IMAGE=fluent/fluent-bit:2.0.9

set -x
docker rm -f flb
mkdir -p /workspaces/data/log
docker run -d --name flb \
-v $PWD/fluent-bit.conf:/fluent-bit/etc/fluent-bit.conf \
-v $PWD/mock.lua:/fluent-bit/etc/mock.lua \
-v /workspaces:/workspaces \
$IMAGE
