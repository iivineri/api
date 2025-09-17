#!/bin/bash

PLATFORM=$(uname)

echo "Platform: $PLATFORM"

CONFIG=$(dirname $0)/develop/docker-compose.yml

NGINX_IP=10.100.10.201

create_domain_name_mapping() {
  echo "Checking virtual host names."

  VHOSTS_FMT="iivineri.local"

  # update /etc/hosts
  if ! grep -q $NGINX_IP /etc/hosts; then
    echo "Updating /etc/hosts."
    sudo sh -c "cat /etc/hosts | grep -v \"$NGINX_IP\" > /tmp/hosts.new"
    sudo sh -c "echo \"$NGINX_IP  $VHOSTS_FMT\" >> /tmp/hosts.new"
    sudo sh -c "cat /tmp/hosts.new > /etc/hosts && rm /tmp/hosts.new"
  fi

  # Check interface for nginx exists
  NGINX_IP_EXISTS=$(ifconfig -a | grep $NGINX_IP)

  #Check if script runs on MacOS device
  IS_MAC=$(uname -a | grep -i "Darwin Kernel Version" | awk '{print $1}')
  #check if it's mac and no nginx IP
  if [[ "$IS_MAC" != "" && "$NGINX_IP_EXISTS" = "" ]]; then
    # Checking Docker version
    echo "Adding $NGINX_IP to lo0 interface"
    sudo sh -c "ifconfig lo0 alias $NGINX_IP"
    NGINX_IP_EXISTS=$(ifconfig -a | grep $NGINX_IP)
  fi

  if [[ "$NGINX_IP_EXISTS" = "" ]]; then
    echo "No Nginx IP. Creating virtual interface."
    CN=$(echo $NGINX_IP | cut -d . -f 4)
    sudo sh -c "ifconfig docker0:ng_$CN $NGINX_IP"
  else
    echo "Nginx IP exists."
  fi

  echo "Domain names mapping finished."
}

export COMPOSE_PROJECT_NAME=iivineri-api
export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1

opts=""
case "$1" in
"build")
  args="--build-arg PLATFORM=$PLATFORM"
  ;;
"up")
  create_domain_name_mapping
  ;;
"start")
  create_domain_name_mapping
  ;;
"shell")
  set -- "exec app /bin/bash"
  ;;
"*") ;;

esac

docker compose --env-file .env --file $CONFIG $opts $@ $args
