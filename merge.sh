#!/bin/bash

echo "" > 00-all-slides.md
for i in */slides*.md 99-end.md; do \
  cat $i >>  00-all-slides.md
  printf "\n\n---\n" >> 00-all-slides.md
done
sed -ie "s;\.\./assets/;assets/;g" "00-all-slides.md"

rm *.mde
