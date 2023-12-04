# CSC482 - Software Deployment - Server Side

## Go [here](https://github.com/hjrose29/agent-hrose3) for the agent(Web Scraping) side!

Skills: GoLang, Docker, AWS(EC2, DynamoDB, IAM, Cloudwatch/trail), HTTP Request Handling/Routing

Server Running [here](http://csc482-alb-1645284752.us-east-1.elb.amazonaws.com/hrose3/status)

<ul>
<li>Endpoints:</li>
<li>'/all' - no parameters, dumps entire table contents</li>
<li>'/rangedSearch' - required parameter: 'ticker' - specifies stock | optional parameters: 'lower' & 'upper' specifies lower and upper UNIX time bounds!</li>
<li>'/search' - required parameters: 'ticker', 'datetime' specifies stock and date time(in UNIX time) respectively</li>
</ul>

Class formed a [microservice](https://en.wikipedia.org/wiki/Microservices) infrastructure!

<b><ins>What I did!</ins></b>

<ol>
<li>Built barebones http server in <ins>GoLang</ins></li>
<li>Fleshed out server with <ins>Gorilla Mux</ins Library </li>
<li>Blocked any non GET requests</li>
<li>Built <ins>middleware</ins> to monitor requests and send them to <ins>Loggly</ins></li>
<li>Containerized server implementation with <ins>multi-stage Docker build</ins> for space efficiency on EC2 instance</li>
<li>Used <ins>DynamoDB query</ins> capabilities to query the database with HTTP requests including optional filter parameters</li>
<li>Added <ins>Data Validation</ins> step to proactively clean input params.</li>
<li>Implemented Continuous Development through <ins>CodeBuild/CodePipeline</ins> AWS's automated deployment services.</li>
<li>Migrated server container to <ins>AWS ECS</ins> and the Application Load Balancer</li>
<li>See it running [here](http://csc482-alb-1645284752.us-east-1.elb.amazonaws.com/hrose3/status)</li>
</ol>

<br>

<b><ins>Building Container</ins>:<b>
docker build . -t <IMAGE_NAME>

<b><ins>Running Container</ins>:<b>
docker run -e <LOGGLY_TOKEN> -p <DOCKER_PORT:OPENED_PORT> <IMAGE_NAME>
