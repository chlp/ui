#!/bin/bash

APP=$(realpath "$1")
ABS_PATH="$(cd "$(dirname "$0")" && pwd)"

cp -f "$ABS_PATH/1_monitor/devices_list_1.json" "$ABS_PATH/1_monitor/devices_list.json"

killall $(basename "$APP") >/dev/null 2>&1 || true

{
  cd "$ABS_PATH/1_monitor"
  $APP &
  pid1=$!
}

{
  cd "$ABS_PATH/2_grpc_device"
  $APP &
  pid2=$!
}

{
  cd "$ABS_PATH/3_rest_device"
  $APP &
  pid3=$!
}

{
  cd "$ABS_PATH/4_rest_fail_device"
  $APP &
  pid4=$!
}

{
  cd "$ABS_PATH/5_rest_warning_device"
  $APP &
  pid5=$!
}

{
  cd "$ABS_PATH/6_rest_device_close_in_30sec"
  $APP &
  pid6=$!
}

trap cleanup SIGINT SIGTERM
cleanup() {
  echo "CLEEAN"
  kill $pid1 $pid2 $pid3 $pid4 $pid5 2>/dev/null || true
  killall $(basename "$APP") >/dev/null 2>&1 || true
  exit 0
}

sleep 30

kill $pid6
echo "Killed 6_rest_device_close_in_30sec"

sleep 5

cp -f "$ABS_PATH/1_monitor/devices_list_2.json" "$ABS_PATH/1_monitor/devices_list.json"
echo "Added 5_rest_warning_device"

sleep 30
cp -f "$ABS_PATH/1_monitor/devices_list_1.json" "$ABS_PATH/1_monitor/devices_list.json"
echo "Removed 5_rest_warning_device"

wait
cleanup
