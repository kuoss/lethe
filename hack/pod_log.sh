namespace="cert-manager"
pod="cert-manager-cainjector-57d984c7c8-lmxjt"
container="cert-manager"

regexPod="Z\[$pod|"
regexContainer="|$container\]\ "

cat /data/log/pod/$namespace/*.log | grep "$regexPod" --color       | tail -10
cat /data/log/pod/$namespace/*.log | grep "$regexContainer" --color | tail -10
