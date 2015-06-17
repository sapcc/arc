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
		createdat timestamp NOT NULL,
		updatedat timestamp NOT NULL
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
		createdat timestamp NOT NULL,
		updatedat timestamp NOT NULL,
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
		createdat timestamp NOT NULL,
		CONSTRAINT log_parts_job_id_number_index UNIQUE (job_id, number)
	)
	WITH (
 	 OIDS=FALSE
	);
	ALTER TABLE logs
		OWNER TO arc;
`
var factsTable = `
	CREATE TABLE IF NOT EXISTS facts
	(
		agent_id text PRIMARY KEY,
		facts jsonb,
		createdat timestamp NOT NULL,
		updatedat timestamp NOT NULL
	)
	WITH (
 	 OIDS=FALSE
	);
	ALTER TABLE facts
		OWNER TO arc;
`

var factsJsonReplaceFunction = `
	CREATE OR REPLACE FUNCTION json_replace(old_data json, new_data json)
	RETURNS json
	IMMUTABLE
	LANGUAGE sql
	AS $$
		SELECT ('{'||string_agg(to_json(key)||':'||value, ',')||'}')::json
		FROM (
			SELECT * FROM json_each(old_data) WHERE key NOT IN (SELECT json_object_keys(new_data))
			UNION ALL
			SELECT * FROM json_each(new_data) WHERE json_typeof(value) <> 'null'
    ) t;
	$$
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

var Tables = [...]string{
	jobsTable,
	logsTable,
	logPartsTable,
	factsTable,
	factsJsonReplaceFunction,
	//agentsTable,
}

// Jobs
var InsertJobQuery = `INSERT INTO jobs(version,sender,requestid,"to",timeout,agent,action,payload,status,createdat,updatedat) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) returning requestid;`
var UpdateJobQuery = `UPDATE jobs SET status=$1,updatedat=$2 WHERE requestid=$3`
var GetAllJobsQuery = "SELECT * FROM jobs order by requestid"
var GetJobQuery = "SELECT * FROM jobs WHERE requestid=$1"
var CleanJobsQuery = `
	UPDATE jobs SET status=3,updatedat=NOW() 
	WHERE requestid IN 
	(
		SELECT DISTINCT requestid
		FROM jobs
		WHERE (createdat <= NOW() - INTERVAL '1 second' * $1 - INTERVAL '1 second' * timeout)
		AND (status=1 OR status=2)
	)
`

// Log
var GetLogQuery = "SELECT content FROM logs WHERE job_id=$1"
var InsertLogQuery = "INSERT INTO logs(job_id,content,createdat,updatedat) VALUES($1,$2,$3,$4) returning job_id"
var UpdateLogQuery = "UPDATE logs SET content=$1,updatedat=$2 WHERE job_id=$3"

// Log parts
var InsertLogPartQuery = `INSERT INTO log_parts(job_id,number,content,final,createdat) VALUES($1,$2,$3,$4,$5) returning job_id;`
var CollectLogPartsQuery = "SELECT array_to_string(array_agg(log_parts.content ORDER BY number, job_id), '') AS content FROM log_parts WHERE job_id=$1"
var DeleteLogPartsQuery = `DELETE FROM log_parts WHERE job_id=$1`

// Facts
var GetAgentsQuery = "SELECT DISTINCT * FROM facts"
var GetAgentQuery = "SELECT agent_id,createdat,updatedat FROM facts WHERE agent_id=$1"
var GetFactQuery = "SELECT facts FROM facts WHERE agent_id=$1"
var InsertFactQuery = `INSERT INTO facts(agent_id,facts,createdat,updatedat) VALUES($1,$2,$3,$4) returning agent_id`
var UpdateFact = `UPDATE facts SET facts=json_replace((SELECT facts::json FROM facts WHERE agent_id=$1),'{"new":"fact"}'::json)::jsonb where agent_id=$1`
