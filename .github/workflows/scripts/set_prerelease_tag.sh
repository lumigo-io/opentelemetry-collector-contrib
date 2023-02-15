# The GitHub Ref is going to be a Pull ref, e.g., `refs/pull/123/head`
TAG="pr-$(echo \"${GITHUB_REF}\" | awk -F/ '{ print $3 }')-${GITHUB_SHA:0:7}-lumigo"

if [[ $TAG =~ ^v[0-9]+\.[0-9]+\.[0-9]+.* ]]
then
    echo "tag=$TAG" >> $GITHUB_OUTPUT
fi
