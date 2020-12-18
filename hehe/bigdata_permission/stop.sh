#!/bin/bash

app_file_name="bigdata_permission"
path=$(dirname $(readlink -f "$0"))
app_path="${path}/${app_file_name}"


# 标题
function title()
{
  echo -n $'\033[33m'
  cat
  echo -n $'\033[0m'
}

# 错误
function error()
{
  echo -n $'\033[31mError: '
  cat
  echo -n $'\033[0m'
  return 1
}>&2


PID=$(ps -def | grep "${app_path}" | grep -v "grep" | awk '{printf $2}')
title <<< "pid: ${PID}"
if [ -n "${PID}" ]; then
    kill -9 ${PID}
    if [ $? -eq 0 ]; then
        title <<< "Program close success."
    else
        error <<< "Program close fail!"
    fi
else
    title <<< "No program to close."
fi

title <<< "done!"