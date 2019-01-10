#!/bin/bash

touch $1.bin
shred -n 1 -s $1 $1.bin
