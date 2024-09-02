#!/bin/sh

source ~/scripts/venv/bin/activate

if [[ ! -d "./generated" ]]; then
    mkdir generated
fi

if [[ ! -d "./generated/aithing-dark" ]]; then
    mkdir generated/aithing-dark
fi
python ~/scripts/svg_to_png.py -f ./aithing-dark.svg -o ./generated/aithing-dark --heights 100 200 300 400 500

if [[ ! -d "./generated/aithing-light" ]]; then
    mkdir generated/aithing-light
fi
python ~/scripts/svg_to_png.py -f ./aithing-light.svg -o ./generated/aithing-light --heights 100 200 300 400 500

if [[ ! -d "./generated/aithing-small" ]]; then
    mkdir generated/aithing-small
fi
python ~/scripts/svg_to_png.py -f ./aithing-small.svg -o ./generated/aithing-small --heights 256

if [[ ! -d "./generated/resume-dark" ]]; then
    mkdir generated/resume-dark
fi
python ~/scripts/svg_to_png.py -f ./resume-dark.svg -o ./generated/resume-dark --heights 100 200 300 400 500

if [[ ! -d "./generated/resume-light" ]]; then
    mkdir generated/resume-light
fi
python ~/scripts/svg_to_png.py -f ./resume-light.svg -o ./generated/resume-light --heights 100 200 300 400 500

if [[ ! -d "./generated/resume-white" ]]; then
    mkdir generated/resume-white
fi
python ~/scripts/svg_to_png.py -f ./resume-white.svg -o ./generated/resume-white --heights 100 200 300 400 500

deactivate