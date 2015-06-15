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
		job_id text NOT NULL,
		content text NOT NULL,
		createdat integer NOT NULL,
		CONSTRAINT index_logs_on_job_id UNIQUE (job_id)
	)
	WITH (
 	 OIDS=FALSE
	);
	ALTER TABLE logs
		OWNER TO arc;
`

var logPartsTable = `
	CREATE TABLE IF NOT EXISTS log_parts
	(
		job_id text NOT NULL,
		number integer NOT NULL,
		content text,
		final boolean,
		createdat integer NOT NULL,
		CONSTRAINT log_parts_job_id_number_index UNIQUE (job_id, number)
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
	logPartsTable,
	//agentsTable,
	//factsTable,
}

// Jobs
var InsertJobQuery = `INSERT INTO jobs(version,sender,requestid,"to",timeout,agent,action,payload,status,createdat,updatedat) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) returning requestid;`
var UpdateJobQuery = `UPDATE jobs SET status=$1,updatedat=$2 WHERE requestid=$3`
var GetAllJobsQuery = "SELECT * FROM jobs order by requestid"
var GetJobQuery = "SELECT * FROM jobs WHERE requestid=$1"

// Log parts
var InsertLogPartQuery = `INSERT INTO log_parts(job_id,number,content,final,createdat) VALUES($1,$2,$3,$4,$5) returning job_id;`
var CollectLogPartsQuery = "SELECT array_to_string(array_agg(log_parts.content ORDER BY number, job_id), '') AS content FROM log_parts WHERE job_id=$1"
