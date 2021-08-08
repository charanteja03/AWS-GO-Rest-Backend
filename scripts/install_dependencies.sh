#!/bin/bash
install_dir="/home/ec2-user/sfbackend"

sudo yum install golang -y

### Check for dir, if not found create it using the mkdir ##
[ ! -d "$dldir" ] && mkdir -p "$install_dir"
