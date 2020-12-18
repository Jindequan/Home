#!/usr/bin/env bash

ENV=$1

EXEC_MAIN_FILE_NAME='bigdata_permission'
EXEC_CLI_FILE_NAME='bigdata_permission_cli'
config_folder='settings'
static_folder='static'
config_name='settings/config_test.json'
config_norm_name='settings/config.json'
template_folder='templates'
start_shell='start.sh'


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

# 解析参数
function parse_param()
{
  echo ${ENV}
  if [[ ${ENV} == 'online' ]]
    then
      config_name='settings/config_online.json'
  fi

  title <<< 'config file::'${config_name}
}



function build()
{
  now_path=`pwd`
  title <<< 'build env:'${ENV}
  title <<< 'build user:'`whoami`
  title <<< ${now_path}
  cd ${now_path}
  title <<< '当前路径：'`pwd`
  go env -w GOPROXY=https://goproxy.cn,direct
  go build --ldflags="-w -s" -o $EXEC_MAIN_FILE_NAME ./main.go
  if [ $? -eq 0 ]; then
    title <<< "go build main_server success"
  else
    error <<< "go build main_server fail!!!"
    return 1
  fi
  title <<< '构建main_server结束'

  go build --ldflags="-w -s" -o $EXEC_CLI_FILE_NAME ./cli.go
  if [ $? -eq 0 ]; then
    title <<< "go build cli success"
  else
    error <<< "go build cli fail!!!"
    return 1
  fi
  title <<< '构建cli结束'
}

function clear_file()
{
  now_path=`pwd`
  title <<< '清理无用文件'
  for file in ${now_path}/*
  do
    echo $file":已删除"
    if [ $file != ${now_path}/${config_folder} -a $file != ${now_path}/${template_folder} -a $file != ${now_path}/${static_folder} -a $file != ${now_path}/$EXEC_MAIN_FILE_NAME -a $file != ${now_path}/$EXEC_CLI_FILE_NAME -a $file != ${now_path}/$start_shell ]; then
      rm -rf $file
    fi
  done
  mv ${now_path}/${config_name} ${now_path}/${config_norm_name}
}


parse_param &&
build &&
clear_file