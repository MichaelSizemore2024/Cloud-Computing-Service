# L3Harris Clinic Team


## Description

This project is a learning project. 
It is designed to be a scalable, load-balancing, modular, and easily distributed Content Delivery Network.
The core aspect of the CDN is the database manager, configured in a Linux environment in Docker.

## Environment Setup

- [ ] Download and open Docker
- [ ] Download and open Visual Studio Code
	- [ ] Terminal: Navigate to cloned Repo
		- [ ] $ docker-compose up -d
- [ ] Containers will download and run
- [ ] To run server/clients see Dev_README.txt

## DBeaver Setup

- [ ] Download and open DBeaver
- [ ] Select Database > Driver Manager
- [ ] Search for Trino
- [ ] Click on the Trino rabbit and select Copy:
	- [ ] Settings: Port 8085 (or whatever port is on Docker trino)
		- [ ] Default User: user
		- [ ] Default Database: cassandra
	- [ ] Libraries: Download/Update button
		- [ ] Select Download after automatic loading
- [ ] Select Database > New Database Connection
	- [ ] Select Trino Copy > Finish

## Add your files

- [ ] [Create](https://docs.gitlab.com/ee/user/project/repository/web_editor.html#create-a-file) or [upload](https://docs.gitlab.com/ee/user/project/repository/web_editor.html#upload-a-file) files
- [ ] [Add files using the command line](https://docs.gitlab.com/ee/gitlab-basics/add-file.html#add-a-file-using-the-command-line) or push an existing Git repository with the following command:

```
cd existing_repo
git remote add origin https://capstone-cs.eng.utah.edu/l3harris-clinic-team/l3harris-clinic-team.git
git branch -M main
git push -uf origin main
```

## Project status

Active Development for University of Utah Spring 2024.
