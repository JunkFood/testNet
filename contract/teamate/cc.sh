#!/bin/bash

docker exec cli peer chaincode install -n teamate -v 0.9 -p github.com/teamate

docker exec cli peer chaincode instantiate -n teamate -v 0.9 -C mychannel -c '{"Args":[]}' -P 'OR("Org1MSP.member")'
sleep 3

docker exec cli peer chaincode invoke -n teamate  -C mychannel -c '{"Args":["registerUser","user1"]}'
sleep 3

docker exec cli peer chaincode invoke -n teamate  -C mychannel -c '{"Args":["joinProject","user1", "project1"]}'
sleep 3

docker exec cli peer chaincode invoke -n teamate  -C mychannel -c '{"Args":["recordScore","user1","project1","5"]}'
sleep 3

docker exec cli peer chaincode invoke -n teamate  -C mychannel -c '{"Args":["readDev","user1"]}'
sleep 3

docker exec cli peer chaincode invoke -n teamate  -C mychannel -c '{"Args":["registerProject","project1","user1"]}'
sleep 3