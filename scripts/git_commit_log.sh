count=$1

git log ${count:+ "-n ${count}"} \
  --pretty=format:'{%n  "commit": "%H",%n  "author": "%aN <%aE>",%n  "date": "%ad",%n  "message": "%f"%n},' |
  perl -pe 'BEGIN{print "["}; END{print "]\n"}' |
  perl -pe 's/},]/}]/'
