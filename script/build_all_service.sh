#!/usr/bin/env bash

source ./style_info.cfg

service_source_root=(
  ../apps/gate/
  ../apps/global/
  ../apps/logic/
  ../apps/login/
  ../apps/management
)

service_names=(
  gate_service
  global_service
  logic_service
  login_service
  cli
)

bin_dir="../bin"
logs_dir="../logs"

if [ ! -d $bin_dir ]; then
  mkdir -p $bin_dir
fi
if [ ! -d $logs_dir ]; then
  mkdir -p $logs_dir
fi

#begin path
begin_path=$PWD
for ((i = 0; i < ${#service_source_root[*]}; i++)); do
  cd $begin_path
  service_path=${service_source_root[$i]}
  cd $service_path && echo -e "${SKY_BLUE_PREFIX}Current directory: $PWD $COLOR_SUFFIX"
  make install && echo -e "${SKY_BLUE_PREFIX}build ${service_names[$i]} success,moving binary file to the bin directory${COLOR_SUFFIX}" &&
    echo -e "${SKY_BLUE_PREFIX}Successful moved ${service_names[$i]} to the bin directory${COLOR_SUFFIX}\n"
done
