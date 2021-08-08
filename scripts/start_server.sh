#!/bin/bash
env_file="/home/ec2-user/sfbackend/.env"

echo "removing old .env file..."
rm "$env_file"
if [ ! -e "$env_file" ]
then
        echo "creating .env file..."
        echo "SERVER_PORT=80"  >> "$env_file"
        echo "BASE_URL=*"  >> "$env_file"
        echo "DISABLE_AUTH=false" >> "$env_file"
else
        echo "no need to create env file"
fi

cd /home/ec2-user/sfbackend
./sfservice >logs.txt 2>errors.txt &
