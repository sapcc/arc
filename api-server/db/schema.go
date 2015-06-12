package db

var jobsTable = `
	CREATE TABLE IF NOT EXISTS jobs
	(
		requestid text PRIMARY KEY,		
		version integer NOT NULL,
		sender text NOT NULL,
		"to" text NOT NULL,
		timeout integer NOT NULL,
		agent text NOT NULL,
		action text NOT NULL,
		payload text NOT NULL,
		status integer NOT NULL,
	)
	WITH (
 	 OIDS=FALSE
	);
	ALTER TABLE jobs
		OWNER TO arc;
`

var logs = `
CREATE TABLE IF NOT EXISTS logs
(
	requestid text PRIMARY KEY,
	id SERIAL,
	payload text NOT NULL
)
WITH (
 OIDS=FALSE
);
ALTER TABLE jobs
	OWNER TO arc;
`

var agentsTable = `
	CREATE TABLE IF NOT EXISTS agents 
	( 
		uid integer PRIMARY KEY,
	)
	WITH (
	 OIDS=FALSE
	);
	ALTER TABLE agents
		OWNER TO arc;
`

var factsTable = `
	CREATE TABLE IF NOT EXISTS facts 
	(
		factID integer PRIMARY KEY NOT null,
		name varchar(255) NOT null,
		value varchar(255) NOT null,
		agent_id int, FOREIGN KEY (agent_id) 
		REFERENCES agents(uid), 
		CONSTRAINT uc_factID UNIQUE (factID, name)
	)
	WITH (
	 OIDS=FALSE
	);
	ALTER TABLE agents
		OWNER TO arc;
`

var Tables = [...]string{
	jobsTable,
	logs,
	agentsTable,
	factsTable,
}

var InsertJobQuery = `INSERT INTO jobs(version,sender,requestid,"to",timeout,agent,action,payload,status) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) returning requestid;`
var UpdateJob = `UPDATE jobs SET status=$1 WHERE requestid=$2`
var GetAllJobsQuery = "SELECT * FROM jobs order by requestid"
var GetJobQuery = "SELECT * FROM jobs WHERE requestid=$1"
