#!/bin/bash

mkdir -p data/json

for f in data/raw/recipes/*; do
    echo ${f##*/}
    ./scripts/troff2json.py $f $* > data/json/${f##*/}.json
done
