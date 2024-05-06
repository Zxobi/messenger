#!/bin/bash
sleep 10

mongosh --host mongo-node1:27017 <<EOF
  var cfg = {
    "_id": "rs0",
    "version": 1,
    "members": [
      {
        "_id": 0,
        "host": "mongo-node1:27017",
        "priority": 2
      },
      {
        "_id": 1,
        "host": "mongo-node2:27017",
        "priority": 0
      },
      {
        "_id": 2,
        "host": "mongo-node3:27017",
        "priority": 0
      }
    ]
  };
  rs.reconfig(cfg, {force: true});
EOF