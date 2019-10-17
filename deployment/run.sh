#!/bin/bash

nohup ./ledmatrix --working-dir=/path/to/ledmatrix > log.txt &
echo $! > pid.txt
