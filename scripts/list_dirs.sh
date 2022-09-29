cd /data/log
a=$(du -s */*)
b=$(du -s --inode */*)
echo "$a" | paste - <(echo "$b") | awk '{printf "%s %s %.1fMi\n", $2, $3-1, $1/1024}' | sed 's|/| |g' | column -t -R 3,4
#awk '{printf "%s %50.1fMi\n", $2, $1/1024}' | column -t
