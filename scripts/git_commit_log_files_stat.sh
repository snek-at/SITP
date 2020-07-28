# cd /Users/nicoschett/Data/GIT/t/intel
count=$1

git log ${count:+ "-n ${count}"} \
    --numstat \
    --format='%H' |
    perl -lawne '
        if (defined $F[1]) {
            print qq#{"insertions": "$F[0]", "deletions": "$F[1]", "path": "$F[2]"},#
        } elsif (defined $F[0]) {
            print qq#],\n"$F[0]": [#
        };
        END{print qq#],#}' |
    tail -n +2 |
    perl -wpe 'BEGIN{print "{"}; END{print "}"}' |
    tr '\n' ' ' |
    perl -wpe 's#(]|}),\s*(]|})#$1$2#g' |
    perl -wpe 's#,\s*?}$#}#'
