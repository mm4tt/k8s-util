#!/bin/bash

set -euo pipefail

git bisect start
git bisect bad  1ad2298
git bisect good 7b62e98

git bisect run ~/run.sh



