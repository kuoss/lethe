cd /data/log; du */*/*.log | awk '{printf "%s %50.1fMi\n", $2, $1/1024}' | sed 's|/| |g' | column -t -R 4
