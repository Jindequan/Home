#!/bin/bash

app_file_name="bigdata_permission"
path=$(dirname $(readlink -f "$0"))

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


# 检查是否有可执行程序
title <<< "Script start..."
app_path="${path}/${app_file_name}"
title <<< "App_path: ${app_path}"

if [ ! -x "${app_path}" ]; then
    error <<< "${app_path} does not exist or no execution permission"
    return
fi

PID=$(ps -def | grep "${app_path}" | grep -v "grep" | awk '{printf $2}')
title <<< "pid: ${PID}"
if [ -n "${PID}" ]; then
    title <<< "Http restart"
    kill -12 ${PID}
    if [ $? -eq 0 ]; then
        title <<< "Restart success! boy!"
    else
        error <<< "Restart fail!"
    fi
else
    title <<< "Http start"
    nohup $app_path > /dev/null 2>&1 &
    if [ $? -eq 0 ]; then
        title <<< "Start success! boy!"
    else
         error <<< "Start fail!"
    fi
fi

title <<< "done!"

