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
		status integer NOT NULL
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
		requestid text PRIMARY KEY,
		id SERIAL,
		payload text NOT NULL
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

var InsertJobQuery = `INSERT INTO jobs(version,sender,requestid,"to",timeout,agent,action,payload,status) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) returning requestid;`
var UpdateJobQuery = `UPDATE jobs SET status=$1 WHERE requestid=$2`
var GetAllJobsQuery = "SELECT * FROM jobs order by requestid"
var GetJobQuery = "SELECT * FROM jobs WHERE requestid=$1"

var InsertLogQuery = `INSERT INTO logs(requestid,payload) VALUES($1,$2) returning requestid;`
var GetLogsQuery = "SELECT * FROM logs WHERE requestid=$1 order by requestid"
