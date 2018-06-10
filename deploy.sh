#!/bin/sh
#########################################################
#
#  Script to pull the latest code and deploy 
#
#
#
#
#########################################################

work_dir=/home/ec2-user
log_file=/home/ec2-user/log.txt


repo=${1}
repo_url=${2}
commit=${3}

echo "Working on ${repo} and commit ${commit}" >> ${log_file}

# Delete existing code to clone it fresh
rm -rf ${work_dir}/${1}


# move to the working directory
cd ${work_dir}

# Clone the repository

git clone ${repo_url}

# change to the repo directory
cd ${work_dir}/${repo}

git fetch --all


git checkout ${commit}


# Perform the maven build for the project
mvn clean install

# Create the docker container
mvn dockerfile:push


# Bring down the app
docker-compose -f  ./src/main/docker/app.yml down --rmi all

# Now run the app
docker-compose -f ./src/main/docker/app.yml up -d 






