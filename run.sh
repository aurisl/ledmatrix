#!/bin/bash

nohup ./ledmatrix > log.txt &
echo $! > pid.txt
