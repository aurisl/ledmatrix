#!/bin/bash

nohup ./led-matrix --working-dir=/path/to/led-matrix > log.txt &
echo $! > pid.txt
