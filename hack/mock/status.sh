set -x
docker ps -a | grep flb$

docker logs flb

find /workspaces/data/log -type f
