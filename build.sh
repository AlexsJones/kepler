#! /bin/bash
#################################################################################
#     File Name           :     build.sh
#     Created By          :     jonesax
#     Creation Date       :     [2017-12-06 15:55]
#     Last Modified       :     [2017-12-06 15:57]
#     Description         :
#################################################################################

go build -i -v -ldflags="-X main.version=$(git describe --always --long --dirty)"
