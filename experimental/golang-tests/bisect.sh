#!/bin/bash

set -euo pipefail

git bisect start
git bisect bad  1ad2298
git bisect good 248444d5eb

git bisect run ~/run.sh



