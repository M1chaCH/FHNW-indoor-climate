-- drop table if exists devices;

create table devices (
    device_id varchar(255) primary key,
    name varchar(255),
    last_ip varchar(16),
    last_reading varchar(255),
    last_reading_time timestamp,
    authorized boolean default false
)