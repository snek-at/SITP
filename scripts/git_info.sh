GIT_URL=$(git config --get remote.origin.url)

JSON_STRING='{"git_url":"'$GIT_URL'"}'

printf $JSON_STRING
