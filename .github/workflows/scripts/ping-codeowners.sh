#!/usr/bin/env bash
#
#   Copyright The OpenTelemetry Authors.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
#


if [ -z "${COMPONENT}"] || [ -z "${ISSUE}" ]; then
    exit 0
fi

result=`grep -c ${COMPONENT} .github/CODEOWNERS`

# there may be more than 1 component matching a label
# if so, try to narrow things down by appending the component
# type to the label
if [[ $result != 1 ]]; then
    COMPONENT_TYPE=`echo ${COMPONENT} | cut -f 1 -d '/'`
    COMPONENT="${COMPONENT}${COMPONENT_TYPE}"
fi

OWNERS=`grep -m 1 ${COMPONENT} .github/CODEOWNERS | sed 's/   */ /g' | cut -f3- -d ' '`

if [ -z "${OWNERS}" ]; then
    exit 0
fi

if [[ "${OWNERS}" =~ "${SENDER}" ]]; then
    echo "Label applied by code owner ${SENDER}"
    exit 0
fi

gh issue comment ${ISSUE} --body "Pinging code owners: ${OWNERS}"

