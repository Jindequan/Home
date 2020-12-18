#!/usr/bin/env bash

declare branch
declare revision=$(printf %x $(date +%s))
PROJECT_NAME='go-basic'
#EXEC_FILE_NAME=$PROJECT_NAME
EXEC_FILE_NAME='go-basic'
TAG_PUSH_FILE_ADDITIONAL='config.json'
START_SHELL='start.sh'
STOP_SHELL='stop.sh'
TEMPLATES_FOLDER='templates'
PUBLIC_FOLDER='public'
config_name = ''

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
  echo ${CI_COMMIT_REF_NAME}
  # 解析构建的branch
  if [[ "${CI_COMMIT_REF_NAME}" == 'master' ]]
  then
    branch='develop'
    config_name='config_sandbox.json'
  elif [[ "${CI_COMMIT_REF_NAME}" == 'pre' ]]
  then
      branch='pre'
      config_name='config_pre.json'
  elif [[ "${CI_COMMIT_REF_NAME}" == 'production' ]]
  then
    branch='production'
    config_name='config_online.json'
  fi
  # 判断参数是否符合条件
  if [[ "${branch}" != 'production' && "${branch}" != 'develop' && "${branch}" != 'pre' ]]
  then
    error <<< 'branch param illegal'
    return
  fi
}

# 初始化git上下文
function init_git_context()
{
  title <<< 'Init GIT Context'
  export CI_REPOSITORY_PUSH_URL=$(echo "${CI_REPOSITORY_URL}" | sed 's/[^@]*/git/' | sed 's/\//:/') &&
  git remote set-url --push origin "${CI_REPOSITORY_PUSH_URL}"
}

# 构建
function build_static()
{
  git_path=`pwd`
  title <<< ${git_path}
  cd ${git_path}
  title <<< '当前路径：'`pwd`
  go env -w GOPROXY=https://goproxy.cn,direct
  go build --ldflags="-w -s" -o $EXEC_FILE_NAME
  if [ $? -eq 0 ]; then
        title <<< "go build success"
    else
        error <<< "go build fail!!!"
        return 1
    fi
  title <<< '构建结束'
}


# 创建标签
function create_tag() {
  title <<< 'Create Tag'
  build_project_path=${git_path}
  for file in $build_project_path/*
  do
      echo $build_project_path/${config_name}
      echo $build_project_path/$EXEC_FILE_NAME
      if [ $file != ${build_project_path}/${config_name} -a $file != ${build_project_path}/$EXEC_FILE_NAME -a $file != ${build_project_path}/$START_SHELL -a $file != ${build_project_path}/$STOP_SHELL -a $file != ${build_project_path}/$TEMPLATES_FOLDER -a $file != ${build_project_path}/$PUBLIC_FOLDER} ]; then
          git rm -f -r $file
      fi
  done
  #重命名配置文件
  mv ${build_project_path}/${config_name} ${build_project_path}/$TAG_PUSH_FILE_ADDITIONAL
  git rm -f -r ${build_project_path}/${config_name}
  git add . &&
  git commit -qm '[GM]GO项目构建' &&
  git tag "${branch}/${revision}" &&
  git push origin "${branch}/${revision}"
}

# 清理标签
function clean_tag {
  local keep=5
  local -a refs
  read -a refs <<< $(git ls-remote --tags "${CI_REPOSITORY_PUSH_URL}" "${branch}"/* | awk '{ print $2}' | sort)
  local total="${#refs[@]}"
  if ! [[ total > keep ]]
  then
    return
  fi
  title <<< 'Clean Tag'
  local i
  local remove=$((total - keep))
  for (( i = 0; i < remove; i ++ ))
  do
    git push origin :"${refs[${i}]}"
  done
}

# 同步执行代码
parse_param && # 解析参数
init_git_context && # 初始化git上下文
build_static && # 构建静态资源
create_tag && # 打上git tag

# 执行代码
clean_tag

