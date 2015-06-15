package db

var jobsTable = `
	CREATE TABLE IF NOT EXISTS jobs
	(
		version integer NOT NULL,
		sender text NOT NULL,
		requestid text PRIMARY KEY,
		"to" text NOT NULL,
		timeout integer NOT NULL,
		agent text NOT NULL,
		action text NOT NULL,
		payload text NOT NULL,
		status integer NOT NULL,
		createdat integer NOT NULL,
		updatedat integer NOT NULL
	)
	WITH (
 	 OIDS=FALSE
	);
	ALTER TABLE jobs
		OWNER TO arc;
`

var logsTable = `
	CREATE TABLE IF NOT EXISTS logs
	(
		requestid text NOT NULL,
		id integer NOT NULL,
		payload text NOT NULL,
		createdat integer NOT NULL,
		CONSTRAINT uc_logID UNIQUE (requestid, id)
	)
	WITH (
 	 OIDS=FALSE
	);
	ALTER TABLE logs
		OWNER TO arc;
`

var agentsTable = `
	CREATE TABLE IF NOT EXISTS agents
	(
		uid text PRIMARY KEY
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
		uid text PRIMARY KEY,
		name varchar(255) NOT null,
		value varchar(255) NOT null
	)
	WITH (
 	 OIDS=FALSE
	);
	ALTER TABLE facts
		OWNER TO arc;
`

var Tables = [...]string{
	jobsTable,
	logsTable,
	//agentsTable,
	//factsTable,
}

var InsertJobQuery = `INSERT INTO jobs(version,sender,requestid,"to",timeout,agent,action,payload,status,createdat,updatedat) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) returning requestid;`
var UpdateJobQuery = `UPDATE jobs SET status=$1,updatedat=$2 WHERE requestid=$3`
var GetAllJobsQuery = "SELECT * FROM jobs order by requestid"
var GetJobQuery = "SELECT * FROM jobs WHERE requestid=$1"

var InsertLogQuery = `INSERT INTO logs(requestid,id,payload,createdat) VALUES($1,$2,$3,$4) returning requestid;`
var GetLogsQuery = "SELECT array_to_string(array_agg(logs.payload ORDER BY id, id), '') AS content FROM logs WHERE requestid=$1"
