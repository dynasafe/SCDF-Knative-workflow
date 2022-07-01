#!/bin/bash

# ============================================parameter setting=========================================

scdf_task_name="java-task03"
ws_url="ws://127.0.0.1:12345/ws"
scdf_server_url=http://dataflow.prd.tanzu

echo "========================================parameter setting========================================"

echo "scdf_task_name: ${scdf_task_name}"
echo "scdf_server_url: ${scdf_server_url}"
echo "ws_url: ${ws_url}"

echo "====================================call scdf to running task===================================="

call_scdf_api_command="curl -k --location --request POST ${scdf_server_url}/tasks/executions?name=${scdf_task_name}"
echo "running command: ${call_scdf_api_command}"
task_id="$(${call_scdf_api_command})"

echo ""
echo "task_id: ${task_id}"

echo "====================================start websocket listening===================================="
python3 -V

python3 <<EOF
from websocket import create_connection

ws = create_connection("${ws_url}")
while True:
    recvContent=ws.recv()
    print("receive task_id: "+recvContent)
    if recvContent=="${task_id}":
        print("task_id: "+recvContent+" is complete")
        break;
ws.close()
EOF

exit 0