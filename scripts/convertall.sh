#!/bin/bash

for f in data/raw/recipes/*; do
    echo ${f##*/}
    ./scripts/troff2json.py -d2 $f > data/json/${f##*/}.json
done
