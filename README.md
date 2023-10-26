## CSC482 - Software Deployment - Server Side

Skills: GoLang, Docker, AWS(EC2, DynamoDB, IAM, Cloudwatch/trail), Loggly

<b><ins>Accomplishments!</ins></b>

<ol>
<li>Built barebones http server in <ins>GoLang</ins></li>
<li>Fleshed out server with <ins>Gorilla Mux</ins Library </li>
<li>Blocked any non GET requests</li>
<li>Built <ins>Middleware</ins> to monitor requests and send them to <ins>Loggly</ins></li>
</ol>


<b><ins>Building Container</ins>:<b>
docker build . -t <IMAGE_NAME>

<b><ins>Running Container</ins>:<b>
docker run -e <LOGGLY_TOKEN> -p <DOCKER_PORT:OPENED_PORT> <IMAGE_NAME>
