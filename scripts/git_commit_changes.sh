version=$1
file=$2
full="$version -- $file"

# Get changes of $file between $version1 and $version2
#
# Example exec: git show 3c46a47a566ef932a463043eeb335137fe4a7f3b -- .github/ISSUE_TEMPLATE/bug_report.md
#
# The git diff --color is needed to prevent git from disabling the color when it is piping
git show --color $full |
  # The grep --color=never is to prevent grep removing the original color and highlighting the matched string.
  # Matching lines that start with red (-) (\e[31m-) or green (+) (\e[32m\+) escape codes
  grep --color=never -E '^\e\[(32m\+|31m-)' |
  # Remove color from output
  sed -E "s/[[:cntrl:]]\[[0-9]{1,3}m//g"
