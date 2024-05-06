CREATE KEYSPACE db_message WITH REPLICATION = { 'class' : 'NetworkTopologyStrategy','DC1' : 3};
use db_message;
CREATE TABLE messages
(
    mid     uuid,
    cid     uuid,
    uid     uuid,
    text    text,
    primary key (cid, mid)
);