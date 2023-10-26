# CSC482 - Software Deployment - Server Side
## Go [here](https://github.com/hjrose29/agent-hrose3) for the agent(Web Scraping) side!

Skills: GoLang, Docker, AWS(EC2, DynamoDB, IAM, Cloudwatch/trail), HTTP Request Handling/Routing

<b><ins>Accomplishments!</ins></b>

<ol>
<li>Built barebones http server in <ins>GoLang</ins></li>
<li>Fleshed out server with <ins>Gorilla Mux</ins Library </li>
<li>Blocked any non GET requests</li>
<li>Built <ins>middleware</ins> to monitor requests and send them to <ins>Loggly</ins></li>
<li>Containerized server implementation with <ins>multi-stage Docker build</ins> for space efficiency on EC2 instance</li>
</ol>

<br>

<b><ins>Building Container</ins>:<b>
docker build . -t <IMAGE_NAME>

<b><ins>Running Container</ins>:<b>
docker run -e <LOGGLY_TOKEN> -p <DOCKER_PORT:OPENED_PORT> <IMAGE_NAME>
