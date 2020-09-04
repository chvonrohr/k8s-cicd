#!/bin/bash

./merge.sh
reveal-md $1 --theme theme/letsboot-white.css --highlight-theme github
